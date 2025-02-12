package service

import (
	"context"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {

	//存起来---存到数据库中---调用下一层rpository
	return svc.repo.Create(ctx, u)
}
