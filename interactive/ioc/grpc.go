package ioc

import (
	grpc2 "github.com/jym0818/webook/interactive/grpc"
	"github.com/jym0818/webook/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(intrServer *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Port      int      `yaml:"port"`
		EtcdAddrs []string `yaml:"etcdAddrs"`
	}

	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	// master 分支
	//err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	intrServer.Register(server)

	return &grpcx.Server{
		Server:    server,
		Port:      cfg.Port,
		EtcdAddrs: cfg.EtcdAddrs,
		Name:      "interactive",
	}
}
