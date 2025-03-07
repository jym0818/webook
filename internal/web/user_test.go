package web

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/service"
	svcmocks "github.com/jym/webook/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) service.UserService
		reqBody string
		//预取响应
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "15904922108@gmail.com", //我期望的参数
					Password: "jy8180900@",
				}).Return(nil) //注册成功返回nil
				return userSvc
			},
			reqBody: `
{
	"email": "15904922108@gmail.com",
	"password": "jy8180900@",
	"confirmPassword": "jy8180900@"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对，bind失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//没有这个，因为参数错误的情况下，就没有预期的调用service的SignuP方法
				//userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
				//	Email:    "15904922108@gmail.com", //我期望的参数
				//	Password: "jy8180900@",
				//}).Return(nil) //注册成功返回nil
				return userSvc
			},
			reqBody: `
{
	"email": "15904922108@gmail.com",
	"password": "jy8180900@",
	"confirmPassword": "jy8180900@",
}
`,
			wantCode: http.StatusBadRequest,
			wantBody: "",
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
				//	Email:    "15904922108@gmail.com", //我期望的参数
				//	Password: "jy8180900@",
				//}).Return(nil) //注册成功返回nil
				return userSvc
			},
			reqBody: `
{
	"email": "15904922108gmail.com",
	"password": "jy8180900@",
	"confirmPassword": "jy8180900@"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱格式不正确",
		},
		{
			name: "密码不正确",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
				//	Email:    "15904922108@gmail.com", //我期望的参数
				//	Password: "jy8180900@",
				//}).Return(nil) //注册成功返回nil
				return userSvc
			},
			reqBody: `
{
	"email": "15904922108@gmail.com",
	"password": "jy8180900",
	"confirmPassword": "jy8180900"
}
`,
			wantCode: http.StatusOK,
			wantBody: "密码必须包含数字、特殊字符，并且长度不能小于 8 位",
		},
		{
			name: "两次密码不相同",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
				//	Email:    "15904922108@gmail.com", //我期望的参数
				//	Password: "jy8180900@",
				//}).Return(nil) //注册成功返回nil
				return userSvc
			},
			reqBody: `
{
	"email": "15904922108@gmail.com",
	"password": "jy8180900@",
	"confirmPassword": "jy8180900@1"
}
`,
			wantCode: http.StatusOK,
			wantBody: "两次密码不同",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "15904922108@gmail.com", //我期望的参数
					Password: "jy8180900@",
				}).Return(service.ErrUserDuplicate) //返回错误
				return userSvc
			},
			reqBody: `
{
	"email": "15904922108@gmail.com",
	"password": "jy8180900@",
	"confirmPassword": "jy8180900@"
}
`,
			wantCode: http.StatusOK,
			wantBody: "重复邮箱，请换一个邮箱",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "15904922108@gmail.com", //我期望的参数
					Password: "jy8180900@",
				}).Return(errors.New("随便返回一个错误")) //其他错误都是系统异常，所以随便返回一个错误
				return userSvc
			},
			reqBody: `
{
	"email": "15904922108@gmail.com",
	"password": "jy8180900@",
	"confirmPassword": "jy8180900@"
}
`,
			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			//用不上codeService
			h := NewUserHandler(tc.mock(ctrl), nil, nil, nil)
			h.RegisterRouters(server)
			//构造请求
			req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
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

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())

		})
	}
}
func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	//我们生成的包中方法-----创建服务
	usersvc := svcmocks.NewMockUserService(ctrl)
	//return 返回什么 要看UserService中的SignUp方法的返回值是什么
	//SignUp的参数也是，不过我们可以通过编译器提示看到
	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
}
