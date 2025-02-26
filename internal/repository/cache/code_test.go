package cache

import (
	"context"
	"errors"
	"github.com/jym/webook/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name  string
		mock  func(ctrl *gomock.Controller) redis.Cmdable
		biz   string
		phone string
		code  string

		wantErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				rd := redismocks.NewMockCmdable(ctrl)

				//看起来是一步调用 其实是两步
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				//redis客户端需要int64转换一下
				cmd.SetVal(int64(0))
				//也就是说cmd的值就是0 错误为nil
				rd.EXPECT().Eval(gomock.Any(), luaSetcode, []string{"phone_code:login:15904922108"}, "989817").Return(cmd)

				return rd
			},
			biz:     "login",
			phone:   "15904922108",
			code:    "989817",
			wantErr: nil,
		},
		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				rd := redismocks.NewMockCmdable(ctrl)

				//看起来是一步调用 其实是两步
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(errors.New("redis error"))
				//redis客户端需要int64转换一下
				//cmd.SetVal(int64(-1))
				//也就是说cmd的值就是0 错误为nil
				rd.EXPECT().Eval(gomock.Any(), luaSetcode, []string{"phone_code:login:15904922108"}, "989817").Return(cmd)

				return rd
			},
			biz:     "login",
			phone:   "15904922108",
			code:    "989817",
			wantErr: errors.New("redis error"),
		},
		{
			name: "发送频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				rd := redismocks.NewMockCmdable(ctrl)

				//看起来是一步调用 其实是两步
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				//redis客户端需要int64转换一下
				cmd.SetVal(int64(-1))
				//也就是说cmd的值就是0 错误为nil
				rd.EXPECT().Eval(gomock.Any(), luaSetcode, []string{"phone_code:login:15904922108"}, "989817").Return(cmd)

				return rd
			},
			biz:     "login",
			phone:   "15904922108",
			code:    "989817",
			wantErr: ErrSetCodeTooMany,
		},
		{
			name: "未知错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				rd := redismocks.NewMockCmdable(ctrl)

				//看起来是一步调用 其实是两步
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				//redis客户端需要int64转换一下
				cmd.SetVal(int64(-2))
				//也就是说cmd的值就是0 错误为nil
				rd.EXPECT().Eval(gomock.Any(), luaSetcode, []string{"phone_code:login:15904922108"}, "989817").Return(cmd)

				return rd
			},
			biz:     "login",
			phone:   "15904922108",
			code:    "989817",
			wantErr: ErrUnknownForCode,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cache := NewCodeCache(tc.mock(ctrl))
			err := cache.Set(context.Background(), tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
