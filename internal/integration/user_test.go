package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/jym/webook/internal/integration/startup"
	"github.com/jym/webook/internal/web"
	"github.com/jym/webook/ioc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_SendLoginSMSCode(t *testing.T) {
	server := startup.InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name string
		//准备数据
		before func(t *testing.T)
		//验证数据--------数据库数据对不都？ redis数据对不对？   所以要你要拿出来
		after func(t *testing.T)

		reqBody  string
		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				//不需要  也就是redis里面什么数据都没有
			},
			after: func(t *testing.T) {
				//验证数据和清理数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

				val, err := rdb.GetDel(ctx, "phone_code:login:15904922108").Result()
				cancel()
				assert.NoError(t, err)
				//验证码是6位
				assert.True(t, len(val) == 6)

			},
			reqBody: `
{
	"phone":"15904922108"
}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{

				Msg: "发送成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				//这个手机号已经有验证码了
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

				_, err := rdb.Set(ctx, "phone_code:login:15904922108", "123456", time.Second*30+time.Minute*9).Result()
				cancel()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				//验证数据和清理数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

				val, err := rdb.GetDel(ctx, "phone_code:login:15904922108").Result()
				cancel()
				assert.NoError(t, err)
				//验证码是6位
				assert.Equal(t, val, "123456")

			},
			reqBody: `
{
	"phone":"15904922108"
}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 4,
				Msg:  "发送频繁",
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

				_, err := rdb.Set(ctx, "phone_code:login:15904922108", "123456", 0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//验证数据和清理数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

				val, err := rdb.GetDel(ctx, "phone_code:login:15904922108").Result()
				cancel()
				assert.NoError(t, err)
				//验证码是6位
				assert.True(t, len(val) == 6)

			},
			reqBody: `
{
	"phone":"15904922108"
}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "手机号码为空",
			before: func(t *testing.T) {
				//不需要  也就是redis里面什么数据都没有
			},
			after: func(t *testing.T) {

			},
			reqBody: `
{
	"phone":""
}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 501001,
				Msg:  "输入有误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)

			//构造请求
			req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			//获取响应
			resp := httptest.NewRecorder()
			//resp.Code  响应码
			//resp.Header()
			//resp.Body

			//这就是HTTP请求进去GIN的入口
			//当你这样调用的时候gin就会处理这个请求，然后响应就会协会resp
			server.ServeHTTP(resp, req)
			var res web.Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, res)
			tc.after(t)

		})
	}
}
