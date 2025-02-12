package repository

import (
	"context"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository/dao"
)

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
func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	//调用底层数据库--->dao层
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindById(id int64) {
	//cache里面找

	//再从dao中找

	//回写cache
}
