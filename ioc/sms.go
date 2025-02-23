package ioc

import (
	"github.com/jym/webook/internal/service/sms"
	"github.com/jym/webook/internal/service/sms/memory"
)

// 可以更改实现  memory或者tencent或者aliyun
func InitSMSService() sms.Service {
	return memory.NewService()
}
