package ioc

import (
	intrv1 "github.com/jym0818/webook/api/proto/gen/intr/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitIntrGRPCClient() intrv1.InteractiveServiceClient {
	type Config struct {
		Addr string
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	cc, err := grpc.Dial(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	res := intrv1.NewInteractiveServiceClient(cc)
	return res
}
