package service

import (
	"context"
	"errors"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicate = repository.ErrUserDuplicate
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

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	//查找
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		//后续接入日志，打印这个错误
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil

}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}
func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {

	//快路径
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// nil会进来-----也就是有用户的
		//其他错误也会进来
		return u, err
	}
	//没有用户 要创建
	//慢路径  一旦服务降级 不走 只保证注册过的用户登录，没注册的用户不提供服务
	//if c.Get("jiangji") == true {
	//	return domain.User{}, err
	//}
	u = domain.User{Phone: phone}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	//没有id啊，怎么办 再找一边
	//这个操作很危险，会遇到主从延迟的问题，如果是主从服务器，那么我们只能让Create接口返回
	return svc.repo.FindByPhone(ctx, phone)

}
