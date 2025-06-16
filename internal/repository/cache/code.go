package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

//go:embed lua/set_code.lua
var luaSetCodeScript string

//go:embed lua/verify_code.lua
var luaVerifyCodeScript string

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("我也不知发生什么了，反正是跟 code 有关")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone string, inputCode string) (bool, error)
}
type codeCache struct {
	cmd redis.Cmdable
}

func (c *codeCache) Set(ctx context.Context, biz, phone, code string) error {
	val, err := c.cmd.Eval(ctx, luaSetCodeScript, []string{c.generate(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch val {
	case 0:
		return nil
	case -2:
		return ErrUnknownForCode
	default:
		return ErrCodeSendTooMany
	}
}

func (c *codeCache) Verify(ctx context.Context, biz, phone string, inputCode string) (bool, error) {
	val, err := c.cmd.Eval(ctx, luaVerifyCodeScript, []string{c.generate(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch val {
	case 0:
		return true, nil
	case -1:
		return false, ErrCodeVerifyTooManyTimes
	default:
		return false, nil
	}
}

func (c *codeCache) generate(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func NewcodeCache(cmd redis.Cmdable) CodeCache {
	return &codeCache{cmd: cmd}
}
