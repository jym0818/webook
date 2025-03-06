package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jym/webook/internal/integration/startup"
	"github.com/jym/webook/internal/repository/dao"
	ijwt "github.com/jym/webook/internal/web/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ArticleTestSuite 测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

// 在所有测试之前  初始化一些内容
func (s *ArticleTestSuite) SetupSuite() {
	s.server = gin.Default()
	s.db = startup.InitDB()
	//模拟用户登录
	s.server.Use(func(c *gin.Context) {
		c.Set("claims", &ijwt.UserClaims{Uid: 123})
	})
	artHdl := startup.InitArticleHandler()
	artHdl.RegisterRouters(s.server)

}

// 结束后运行
func (s *ArticleTestSuite) TearDownSuite() {
	//清空数据库 将自增主键回复为1
	s.db.Exec("TRUNCATE TABLE articles")
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

// 使用了测试套件后 不需要在测试方法中添加t *testing.T
func (s *ArticleTestSuite) TestABC() {
	s.T().Log("hello 这是测试套件")
}
func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name string
		//预期输入
		art Article
		//http响应码
		wantCode int
		//预期响应
		wantRes Result[int64]

		before func(t *testing.T)
		after  func(t *testing.T)
	}{
		{
			name: "新建帖子---保存成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.Article
				err := s.db.Where("id=?", 1).First(&art).Error

				assert.NoError(t, err)
				//没办法判断Ctime和Utime，所以比较一下，设为0
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
				}, art)

			},
			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "ok",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			//构造请求
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/articles/edit", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			//获取响应
			resp := httptest.NewRecorder()
			//resp.Code  响应码
			//resp.Header()
			//resp.Body

			//这就是HTTP请求进去GIN的入口
			//当你这样调用的时候gin就会处理这个请求，然后响应就会协会resp
			s.server.ServeHTTP(resp, req)
			var res Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantRes, res)
			tc.after(t)

		})
	}
}

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
