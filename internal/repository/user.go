package repository

import (
	"context"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserReposity(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
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

func (r *UserRepository) FindById(id int64) {
	//cache里面找

	//再从dao中找

	//回写cache
}
