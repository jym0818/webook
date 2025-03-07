package web

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jym/webook/internal/domain"
	"github.com/jym/webook/internal/service"
	svcmocks "github.com/jym/webook/internal/service/mocks"
	ijwt "github.com/jym/webook/internal/web/jwt"
	"github.com/jym/webook/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) service.ArticleService
		reqBody string
		//预取响应
		wantCode int
		wantRes  Result
	}{
		{
			name: "新建并发表",
			reqBody: `
{
	"title":"我的标题",
	"content": "我的内容"
}
`,
			wantCode: http.StatusOK,
			wantRes: Result{
				//强制指定为float64
				Data: float64(1),
				Msg:  "ok",
			},
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			server.Use(func(c *gin.Context) {
				c.Set("claims", &ijwt.UserClaims{Uid: 123})
			})
			h := NewArticleHandler(tc.mock(ctrl), logger.NewNoOpLogger())
			h.RegisterRouters(server)

			req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/articles/publish", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			//获取响应
			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			var webRes Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)

		})
	}
}

//在json里面字段是any  数字类型是float64
//如果还是json，会转换为map
