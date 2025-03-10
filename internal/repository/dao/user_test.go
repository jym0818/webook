package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name string
		//这里是sqlmock，不是gomock
		mock    func(t *testing.T) *sql.DB
		ctx     context.Context
		user    User
		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				//这个是预期的正则表达式
				mock.ExpectExec("INSERT INTO `users`.*").WillReturnResult(res)
				require.NoError(t, err)
				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "test@test.com",
					Valid:  true,
				},
				Password: "123456",
			},

			wantErr: nil,
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()

				//这个是预期的正则表达式
				mock.ExpectExec("INSERT INTO `users`.*").WillReturnError(&mysql.MySQLError{
					Number: 1062,
				})
				require.NoError(t, err)
				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "test@test.com",
					Valid:  true,
				},
				Password: "123456",
			},

			wantErr: ErrUserDuplicate,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()

				//这个是预期的正则表达式
				mock.ExpectExec("INSERT INTO `users`.*").WillReturnError(errors.New("db error"))
				require.NoError(t, err)
				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "test@test.com",
					Valid:  true,
				},
				Password: "123456",
			},

			wantErr: errors.New("db error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name+" insert", func(t *testing.T) {

			db, err := gorm.Open(gormMysql.New(gormMysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})
			//初始化db不能出错，我们要断言
			assert.NoError(t, err)
			d := NewUserDAO(db)
			err = d.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)

		})
	}
}
