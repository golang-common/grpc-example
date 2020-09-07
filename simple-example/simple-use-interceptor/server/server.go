// @Author: Perry
// @Date  : 2020/8/26
// @Desc  : 带拦截器的简单服务器

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"grpc-example/simple-example/simple-proto-out/hello"
	"log"
	"net"
	"time"
)

const serverAddr = "localhost:8080"

/*方法实现*/
type HelloServer struct{}

func (*HelloServer) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	fmt.Println("function called,client message received")
	lastUpdate := new(timestamp.Timestamp)
	lastUpdate.Seconds = time.Now().Unix()
	return &hello.HelloReply{Resp: "this is a reply message from server,client message", LastUpdate: lastUpdate}, nil
}

// 打印客户端发来context中包含的元数据信息
func printContextMeta(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("context中不包含任何信息")
	}
	for k, v := range md {
		fmt.Printf("key=%s , value=%s\n", k, v)
	}
	return nil
}

// 服务端拦截器处理函数,需符合grpc.UnaryServerInterceptor类型定义
// 本拦截器打印context中包含的元数据信息，以及各请求参数信息
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	/*
		ctx 服务端接收的上下文信息
		req 客户端发起的请求信息，具体类型为protoc中定义的消息结构体
		info 客户端请求对应的具体protoc中的方法等信息，本例中为SayHello
		handler 处理方法，调用handler方法即为拦截后的处理，在此之前为拦截前的处理
		返回：resp,err 为调用handler之后的返回结果
	*/
	err = printContextMeta(ctx)
	fmt.Printf("req = %+v\n", req.(*hello.HelloRequest).Req)
	fmt.Printf("info.server = %+v\n", info.Server)
	fmt.Printf("info.method = %+v\n", info.FullMethod)
	m, err := handler(ctx, req)
	if err != nil {
		fmt.Println(err.Error())
	}
	return m, err
}

func main() {
	// 创建网络连接
	conn, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// 定义拦截器
	// 新建grpc服务,增加grpc选项
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(unaryInterceptor))

	hello.RegisterHelloServer(server, &HelloServer{})
	reflection.Register(server)
	fmt.Println("server started")
	if err = server.Serve(conn); err != nil {
		log.Fatal(err.Error())
	}
}