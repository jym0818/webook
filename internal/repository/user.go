package repository

import (
	"context"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

type userRepository struct {
	dao dao.UserDAO
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{Id: user.Id, Email: user.Email, Password: user.Password}, nil
}

func NewuserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{dao: dao}
}

func (repo *userRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, dao.User{Email: user.Email, Password: user.Password})
}
