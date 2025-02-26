package ioc

import "github.com/jym/webook/internal/service/oauth2/wechat"

func InitOAuth2WechatService() wechat.Service {
	appId := "123456789"
	appSecret := "123456789"
	return wechat.Newservice(appId, appSecret)
}
