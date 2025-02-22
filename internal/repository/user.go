package repository

import (
	"context"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository/cache"
	"github.com/jym/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserReposity(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

// 命名为Create 因为在这一层级repository中已经没有signup的概念了
// 数据传递通常为结构体，而不是结构体指针
func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	//调用底层数据库--->dao层
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	//调用底层数据库--->dao层
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		//如果出错，返回空结构体
		return domain.User{}, err
	}
	return domain.User{Id: u.Id, Email: u.Email, Password: u.Password}, nil

}

func (repo *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	//cache里面找
	user, err := repo.cache.Get(ctx, id)
	//err几种情况
	//1.缓存有数据-----err为nil
	//2.缓存没有数据 -----err为
	//3.缓存出错------err为系统错误，直接返回
	if err == nil {
		return user, nil
	}
	//err为其他错误（系统错误），怎么办？要不要去数据库加载？
	//如果现在redis崩溃了（缓存雪崩、穿透了），我们如果让这些请求去数据库上加载，数据库不就崩了吗
	//选加载------万一redis真崩了，我们必须保护住我们的数据库
	//选不加载-----用户体验差一点
	//选加载，我们方案是数据库限流;用orm的middleware
	ue, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u := domain.User{Id: ue.Id, Email: ue.Email, Password: ue.Password}
	//回写cache
	err = repo.cache.Set(ctx, u)
	if err != nil {
		//缓存设置失败，我这里怎么办，要不要返回err
		//不需要返回，打个日志就可以了,要监控好，防止redis崩了
	}
	return u, err
}
