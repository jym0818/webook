package service

import (
	"context"
	"fmt"
	"github.com/jym0818/webook/internal/repository"
	"github.com/jym0818/webook/internal/service/sms"
	"math/rand"
)

var (
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrUnknownForCode         = repository.ErrUnknownForCode
)

const tpl = "123456"

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz string, phone string, code string) (bool, error)
}

type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func (svc *codeService) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, code)
}

func (svc *codeService) Send(ctx context.Context, biz, phone string) error {
	//生成验证码
	code := svc.generate()
	//redis存储
	err := svc.repo.Store(ctx, code, biz, phone)
	if err != nil {
		return err
	}
	//发送
	err = svc.smsSvc.Send(ctx, tpl, []string{code}, phone)

	if err != nil {
		//记录日志就可以了
	}
	return nil
}

func (svc *codeService) generate() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func NewcodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}
