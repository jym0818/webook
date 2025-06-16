package memory

import (
	"context"
	"fmt"
	"github.com/jym0818/webook/internal/service/sms"
)

type Service struct{}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}

func NewService() sms.Service {
	return &Service{}
}
