package main


import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	pb "ProxyPool/proto"

)

func main() {
	//ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "0.0.0.0:8082", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}
	//time.Sleep(11 * time.Second)
	//cancel()

	client := pb.NewProxyPoolClient(conn)
	resp, err := client.CallProxy(ctx, &pb.Req{})
	if err != nil {
		log.Fatal(err)
	}
	logrus.Println(resp.Proxy)
	logrus.Println(resp.Message)
}
