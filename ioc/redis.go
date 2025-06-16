package ioc

import "github.com/redis/go-redis/v9"

func InitRedis() redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr: "118.25.44.1:6379",
	})
	return cmd
}
