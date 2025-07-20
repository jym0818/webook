package grpcx

import (
	"context"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
)

type Server struct {
	*grpc.Server
	Port      int
	EtcdAddrs []string
	Name      string
	kaCancel  func()
	em        endpoints.Manager
	client    *etcdv3.Client
	key       string
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		return err
	}
	err = s.register()
	if err != nil {
		return err
	}
	// 就是直接启动，我现在要嵌入服务注册过程
	// 这边会阻塞，类似与 gin.Run
	return s.Server.Serve(l)
}

func (s *Server) register() error {
	client, err := etcdv3.New(etcdv3.Config{
		Endpoints: s.EtcdAddrs,
	})
	if err != nil {
		return err
	}
	s.client = client
	// endpoint 以服务为维度。一个服务一个 Manager
	em, err := endpoints.NewManager(client, "service/"+s.Name)
	if err != nil {
		return err
	}
	addr := netx.GetOutboundIP() + ":" + strconv.Itoa(s.Port)
	key := "service/" + s.Name + "/" + addr
	s.key = key
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 你可以做成配置的
	var ttl int64 = 30
	leaseResp, err := client.Grant(ctx, ttl)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, etcdv3.WithLease(leaseResp.ID))

	kaCtx, kaCancel := context.WithCancel(context.Background())
	s.kaCancel = kaCancel
	ch, err := client.KeepAlive(kaCtx, leaseResp.ID)
	if err != nil {
		return err
	}
	go func() {
		for kaResp := range ch {
			// 正常就是打印一下 DEBUG 日志啥的
			zap.L().Debug(kaResp.String())
		}
	}()
	return nil
}

// Close 你可以叫做 Shutdown
func (s *Server) Close() error {
	if s.kaCancel != nil {
		s.kaCancel()
	}
	if s.em != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := s.em.DeleteEndpoint(ctx, s.key)
		if err != nil {
			return err
		}
	}
	if s.client != nil {
		err := s.client.Close()
		if err != nil {
			return err
		}
	}
	s.Server.GracefulStop()
	return nil
}
