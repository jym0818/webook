package repository

import (
	"context"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
}

type userRepository struct {
	dao dao.UserDAO
}

func NewuserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{dao: dao}
}

func (repo *userRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, dao.User{Email: user.Email, Password: user.Password})
}
