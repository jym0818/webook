package ioc

import (
	"github.com/jym0818/webook/internal/service/sms"
	"github.com/jym0818/webook/internal/service/sms/memory"
)

func InitSMS() sms.Service {
	smsSvc := memory.NewService()
	return smsSvc
}
