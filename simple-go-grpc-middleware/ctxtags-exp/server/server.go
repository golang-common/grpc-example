// @Author: Perry
// @Date  : 2020/8/26
// @Desc  :

package main

import (
	"context"
	"fmt"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-example/simple-proto-out/hello"
	"log"
	"net"
	"time"
)

const serverAddr = "localhost:8080"

/*方法实现*/
type HelloServer struct{}

func (*HelloServer) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	fmt.Println("function called,client message=\n", req.Req)
	lastUpdate := new(timestamp.Timestamp)
	lastUpdate.Seconds = time.Now().Unix()
	return &hello.HelloReply{Resp: "this is a reply message from server,client message", LastUpdate: lastUpdate}, nil
}

func main() {
	conn, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 定义拦截器option
	opts := []grpc_ctxtags.Option{
		grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.TagBasedRequestFieldExtractor("log_fields")),
	}
	// 创建server时增加相关拦截器选项
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_ctxtags.StreamServerInterceptor(opts...)),
		grpc.UnaryInterceptor(grpc_ctxtags.UnaryServerInterceptor(opts...)),
	)
	hello.RegisterHelloServer(server, &HelloServer{})
	reflection.Register(server)
	fmt.Println("server started")
	if err = server.Serve(conn); err != nil {
		log.Fatal(err.Error())
	}
}
