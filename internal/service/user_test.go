package service

import (
	"context"
	"errors"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/repository"
	repomocks "github.com/jym/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func Test_userService_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository
		//输入
		email    string
		password string

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@test.com").Return(domain.User{
					Email:    "test@test.com",
					Phone:    "15904922108",
					Password: "$2a$10$AXLy/R7vjXa5ziWfY21W.ORurXEKRSYvGjMZQCkDvz3SENEUlNmk6", //密码加密：
					Ctime:    now,
				}, nil)
				return repo
			},
			email:    "test@test.com",
			password: "jy1206125811@",
			wantUser: domain.User{
				Email:    "test@test.com",
				Phone:    "15904922108",
				Password: "$2a$10$AXLy/R7vjXa5ziWfY21W.ORurXEKRSYvGjMZQCkDvz3SENEUlNmk6",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@test.com").Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "test@test.com",
			password: "jy1206125811@",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@test.com").Return(domain.User{}, errors.New("系统错误"))
				return repo
			},
			email:    "test@test.com",
			password: "jy1206125811@",
			wantUser: domain.User{},
			wantErr:  errors.New("系统错误"),
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@test.com").Return(domain.User{
					Email:    "test@test.com",
					Phone:    "15904922108",
					Password: "$2a$10$AXLy/R7vjXa5ziWfY21W.ORurXEKRSYvGjMZQCkDvz3SENEUlNmk6", //密码加密：
					Ctime:    now,
				}, nil)
				return repo
			},
			email:    "test@test.com",
			password: "1jy1206125811@", //随便输出一个密码
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			//具体测试代码
			svc := NewUserService(tc.mock(ctrl), nil)
			user, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("jy1206125811@"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
