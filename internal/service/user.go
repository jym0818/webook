package service

import (
	"context"
	"errors"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号或者密码错误")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	//加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	u.Password = string(hash)
	//存起来---存到数据库中---调用下一层rpository
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email, password string) error {
	//查找
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return ErrInvalidUserOrPassword
	}
	if err != nil {
		return err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		//后续接入日志，打印这个错误
		return ErrInvalidUserOrPassword
	}
	return nil

}
