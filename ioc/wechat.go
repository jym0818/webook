package ioc

import (
	"github.com/jym0818/webook/internal/service/oauth2/wechat"
	"github.com/jym0818/webook/internal/web"
)

func InitWechat() wechat.Service {
	service := wechat.Newservice("123456", "123456")
	return service
}

func InitWechatCfg() web.Config {
	return web.Config{Secure: false}
}
