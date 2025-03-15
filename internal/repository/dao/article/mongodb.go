package article

import (
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBDAO struct {
	//实际上client和database用不到
	//client *mongo.Client
	////代表webook
	//database *mongo.Database
	//代表制作库
	col *mongo.Collection
	//代表线上库
	liveCol *mongo.Collection
	node    *snowflake.Node
}

func (m *MongoDBDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	//生成
	id := m.node.Generate().Int64()
	art.Id = id
	_, err := m.col.InsertOne(ctx, art)
	//返回的是什么？这不是自增主键  而是一个[12]byte 自动生成的唯一标识
	//return res.InsertedID, nil
	return id, err
}

func (m *MongoDBDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	res, err := m.col.UpdateOne(ctx, filter, bson.D{bson.E{"$set", bson.M{
		"context": art.Content,
		"title":   art.Title,
		"utime":   now,
		"status":  art.Status,
	}}})

	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("可能出现用户修改其他人，受到攻击了？")
	}
	return nil

}

func (m *MongoDBDAO) Sync(ctx context.Context, art Article) (int64, error) {
	//没办法使用事务
	//1.保存制作库
	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	//2.线上库 upsert
	now := time.Now().UnixMilli()
	art.Utime = now
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	update := bson.E{"$set", art}
	upsert := bson.E{"$setOnInsert", bson.D{bson.E{"ctime", now}}}
	_, err = m.liveCol.UpdateOne(ctx, filter,
		bson.D{update, upsert}, options.Update().SetUpsert(true))
	return id, err

}

func (m *MongoDBDAO) Upsert(ctx context.Context, art PublishArticle) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) SyncStatus(ctx context.Context, id int64, uid int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}

func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) ArticleDAO {
	return &MongoDBDAO{
		node:    node,
		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
	}
}
func InitCollections(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	index := []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "author_id", Value: 1},
				bson.E{Key: "ctime", Value: 1},
			},
			Options: options.Index(),
		},
	}
	_, err := db.Collection("articles").Indexes().
		CreateMany(ctx, index)
	if err != nil {
		return err
	}
	_, err = db.Collection("published_articles").Indexes().
		CreateMany(ctx, index)
	return err
}
