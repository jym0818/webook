package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository/cache"
	cachemocks "github.com/jym/webook/internal/repository/cache/mocks"
	"github.com/jym/webook/internal/repository/dao"
	daomocks "github.com/jym/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	now := time.Now()
	//now包含纳秒 你必须去掉纳秒
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		id   int64

		wantUser domain.User
		wantErr  error
	}{
		//第一个用例  一般是路径最长的，方便测试和复制
		{
			name: "缓存未命中,但是查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				//缓存未命中  查了缓存没结果  又查了数据库才有结果
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDAO(ctrl)
				//因为是接口类型，直接输入1会被定义为int类型，所以指定int64，这里的id必须与参数的id相同
				c.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, errors.New("未查到或者系统错误"))
				//这里的id必须与参数的id相同 dao.User应该与wantUser相同
				d.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{
					Id:       1,
					Email:    sql.NullString{String: "test@test.com", Valid: true},
					Password: "xxxxxx",
					Phone:    sql.NullString{String: "13940652390", Valid: true},
					Ctime:    now.UnixMilli(),
				}, nil)
				//domain.User与上面相同
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       1,
					Email:    "test@test.com",
					Password: "xxxxxx",
					Phone:    "13940652390",
					Ctime:    now,
				}).Return(nil)
				return d, c
			},
			id: 1,

			wantUser: domain.User{
				Id:       1,
				Email:    "test@test.com",
				Password: "xxxxxx",
				Phone:    "13940652390",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				//缓存未命中  查了缓存没结果  又查了数据库才有结果
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDAO(ctrl)
				//因为是接口类型，直接输入1会被定义为int类型，所以指定int64，这里的id必须与参数的id相同
				c.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{
					Id:       1,
					Email:    "test@test.com",
					Password: "xxxxxx",
					Phone:    "13940652390",
					Ctime:    now,
				}, nil)
				return d, c
			},
			id: 1,

			wantUser: domain.User{
				Id:       1,
				Email:    "test@test.com",
				Password: "xxxxxx",
				Phone:    "13940652390",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "数据库查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				//缓存未命中  查了缓存没结果  又查了数据库才有结果
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDAO(ctrl)
				//因为是接口类型，直接输入1会被定义为int类型，所以指定int64，这里的id必须与参数的id相同
				c.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, errors.New("未查到或者系统错误"))
				//这里的id必须与参数的id相同 dao.User应该与wantUser相同
				d.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{}, errors.New("db错误"))

				return d, c
			},
			id: 1,

			wantUser: domain.User{},
			wantErr:  errors.New("db错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			//1.需要两个参数，去搞参数吧
			//mock来获取呗
			//创建对象
			repo := NewUserReposity(ud, uc)
			//测试的是repository.FindById方法
			// 两个参数
			res, err := repo.FindById(context.Background(), tc.id)
			//一般先断言error
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, res)

		})
	}
}
