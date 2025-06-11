package service

import (
	"context"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/repository"
)

type UserService interface {
	Signup(ctx context.Context, user domain.User) error
}

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail

type userService struct {
	repo repository.UserRepository
}

func NewuserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx context.Context, user domain.User) error {
	return svc.repo.Create(ctx, user)
}
