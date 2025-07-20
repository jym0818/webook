package ioc

import (
	intrv1 "github.com/jym0818/webook/api/proto/gen/intr/v1"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitEtcd() *clientv3.Client {
	var cfg clientv3.Config
	err := viper.UnmarshalKey("etcd", &cfg)
	if err != nil {
		panic(err)
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return cli
}

func InitIntrGRPCClient(client *clientv3.Client) intrv1.InteractiveServiceClient {
	type Config struct {
		Name string
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}

	bd, err := resolver.NewBuilder(client)
	if err != nil {
		panic(err)
	}
	cc, err := grpc.Dial("etcd:///service/"+cfg.Name, grpc.WithResolvers(bd), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	res := intrv1.NewInteractiveServiceClient(cc)
	return res
}
