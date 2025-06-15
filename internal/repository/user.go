package repository

import (
	"context"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/repository/cache"
	"github.com/jym0818/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, uid int64) (domain.User, error)
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
	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}
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
	return domain.User{Id: user.Id, Email: user.Email, Password: user.Password}, nil
}

func NewuserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &userRepository{dao: dao, cache: cache}
}

func (repo *userRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, dao.User{Email: user.Email, Password: user.Password})
}
