package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ErrSetCodeTooMany = errors.New("发送太频繁")
var ErrCodeVerifyTooManyTimes = errors.New("验证次数太多了")
var ErrUnknownForCode = errors.New("发送验证码遇到未知错误")

// 它通过//go:embed 指令，可以在编译阶段将静态资源文件打包进编译好的程序中，并提供访问这些文件的能力
//编译器会在编译的时候，把set_code代码放进这个luaSetCode变量里面

//go:embed lua/set_code.lua
var luaSetcode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}
type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}
func (c *RedisCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {

	res, err := c.client.Eval(ctx, luaSetcode, []string{c.Key(biz, phone)}, code).Int()

	if err != nil {
		return err
	}
	switch res {
	case 0:
		//完全正常
		return nil
	case -1:
		//发送频繁，返回特定错误

		zap.L().Warn("短信发送太频繁")
		//对应的告警系统要配置规则
		//比如1分钟内出现100次短信发送太频繁就告警，这意味着有人在搞你
		//你要去看看，是不是真的有人在搞，你应该去触发一些安全策略
		return ErrSetCodeTooMany
	default:
		//系统错误
		return ErrUnknownForCode
	}
}

func (c *RedisCodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	//获取
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.Key(biz, phone)}, code).Int()
	fmt.Println(res)
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		//正常来说，如果频繁出现这个错误，应该告警
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
	default:
		return false, nil
	}
}
