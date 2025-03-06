package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jym/webook/internal/domain"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

var redirectURI = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

type service struct {
	appId     string
	appSecret string
	client    *http.Client
}

func Newservice(appId string, appSecret string) Service {
	//为什么client没有依赖注入   因为除了这里别的地方没有使用 所以偷懒了
	return &service{appId: appId, appSecret: appSecret, client: http.DefaultClient}
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redire"
	//随机生成的

	return fmt.Sprintf(urlPattern, s.appId, redirectURI, state), nil
}

func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, s.appId, redirectURI, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("微信返回错误信息:%s", res.ErrMsg)
	}

	//日志记录  openId和UnionId不是敏感信息
	zap.L().Info("调用微信，拿到用户信息", zap.String("unionId", res.UnionId), zap.String("openId", res.OpenId))

	return domain.WechatInfo{
		OpenID:  res.OpenId,
		UnionID: res.UnionId,
	}, nil

}

// 根据腾讯文档的返回数据定义的结构体
type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	UnionId string `json:"unionid"`
	OpenId  string `json:"openid"`
}
