package service

import (
	"context"
	"fmt"
	"github.com/jym/webook/internal/repository"
	"github.com/jym/webook/internal/service/sms"
	"math/rand"
)

const codeTplId = "2367159"

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{repo: repo, smsSvc: smsSvc}
}

// 发送验证码
// 多个业务会发送验证码 biz来区分业务场景
func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	//生成验证码
	code := svc.generateCode()
	//放进redis
	err := svc.repo.Store(ctx, biz, code, phone)

	if err != nil {
		return err
	}
	//发送出去
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	//if err != nil {
	//	// 如果前面成功了，发送这步骤失败了怎么办
	//	//这意味这redis有这个验证码，但是用户没有收到
	//	//我需不需要会redis删除验证码
	//	//有可能是超时的err，你根本无法判断是否发出
	//	//可以接受发送失败，忽略
	//}
	return err

}

//func (svc *CodeService) Verfiy(ctx context.Context, biz string, inputCode string, phone string) (bool, error) {
//
//}

func (svc *CodeService) generateCode() string {
	//0 - 999999 之间 包含0 和999999
	num := rand.Intn(1000000)
	//不够6位 加上前缀0
	return fmt.Sprintf("%6d", num)
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}
