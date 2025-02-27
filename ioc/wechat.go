package ioc

import (
	"github.com/jym/webook/internal/service/oauth2/wechat"
	"github.com/jym/webook/internal/web"
)

func InitOAuth2WechatService() wechat.Service {
	appId := "123456789"
	appSecret := "123456789"
	return wechat.Newservice(appId, appSecret)
}

func NewWechatHandler() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
