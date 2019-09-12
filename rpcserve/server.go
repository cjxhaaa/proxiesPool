package rpcserve

import (
	"ProxyPool/pkg/crawler"
	pb "ProxyPool/proto"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)


type Service struct {
	Port  string
	Pool  *crawler.Pool
}

func (s *Service) CallProxy(ctx context.Context,req *pb.Req) (*pb.Resp, error) {
	proxy := s.Pool.GetOneProxy()
	if proxy != "" {
		return &pb.Resp{Proxy:proxy, Message:pb.Resp_SUCCESS}, nil
	}
	return &pb.Resp{Message:pb.Resp_FAILED}, nil
}

func (s *Service)RunGrpcServer() {
	lis, err := net.Listen("tcp", "0.0.0.0:"+s.Port)
	if err != nil {
		panic(err)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		logrus.Infof("")
		logrus.Infof("reveive signal '%v'", s)
		os.Exit(1)
	}()
	sg := grpc.NewServer()
	pb.RegisterProxyPoolServer(sg, s)
	sg.Serve(lis)
}

