package repository

import (
	"context"
	"database/sql"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/repository/cache"
	"github.com/jym0818/webook/internal/repository/dao"
	"time"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, uid int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByWechat(ctx context.Context, openId string) (domain.User, error)
}

type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func (repo *userRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	u, err := repo.cache.Get(ctx, uid)
	if err == nil {
		return u, nil
	}

	ue, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	u = repo.toDomain(ue)
	//回写缓存
	go func() {
		err = repo.cache.Set(ctx, u)
		if err != nil {
			//记录日志  prometheus
		}
	}()
	return u, nil
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(user), nil
}

func (repo *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(user), nil
}

func NewuserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &userRepository{dao: dao, cache: cache}
}

func (repo *userRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(user))
}

func (repo *userRepository) toEntity(user domain.User) dao.User {
	return dao.User{
		Id:       user.Id,
		Ctime:    user.Ctime.UnixMilli(),
		Utime:    user.Utime.UnixMilli(),
		Password: user.Password,
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		WechatOpenID: sql.NullString{
			String: user.WechatInfo.OpenID,
			Valid:  user.WechatInfo.OpenID != "",
		},
		WechatUnionID: sql.NullString{
			String: user.WechatInfo.UnionID,
			Valid:  user.WechatInfo.UnionID != "",
		},
		Nickname: user.Nickname,
	}
}

func (repo *userRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Utime:    time.UnixMilli(user.Utime),
		Ctime:    time.UnixMilli(user.Ctime),
		Phone:    user.Phone.String,
		Password: user.Password,
		Email:    user.Email.String,
		WechatInfo: domain.WechatInfo{
			OpenID:  user.WechatOpenID.String,
			UnionID: user.WechatUnionID.String,
		},
		Nickname: user.Nickname,
	}
}

func (repo *userRepository) FindByWechat(ctx context.Context, openId string) (domain.User, error) {
	user, err := repo.dao.FindByWechat(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(user), nil
}
