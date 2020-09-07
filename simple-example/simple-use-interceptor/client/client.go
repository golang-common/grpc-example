// @Author: Perry
// @Date  : 2020/8/26
// @Desc  : 带拦截器的简单客户端

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"grpc-example/simple-example/simple-proto-out/hello"
	"log"
)

const serverAddr = "localhost:8080"

// 客户端拦截器，在客户端发送前需要统一完成的工作
func unaryInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	/*
		ctx 客户端发送时的上下文信息，此处可以插入相关的键值信息
		method 调用的方法名称

	*/
	var md = make(map[string][]string)
	md["appid"] = []string{"123456"}
	md["name"] = []string{"daipengyuan", "lyonsdpy"}
	newCtx := metadata.NewOutgoingContext(ctx, md)

	// 在调用之前打印reply，应该为空
	fmt.Println("resp = ", reply.(*hello.HelloReply).Resp)

	err := invoker(newCtx, method, req, reply, cc, opts...) //调用方法继续执行，之前为调用之前，之后为调用之后

	// 打印context中的元数据信息
	if mdNew, ok := metadata.FromOutgoingContext(newCtx); ok {
		for k, v := range mdNew {
			fmt.Println("key = ", k)
			fmt.Println("val = ", v)
		}
	}
	// 打印method
	fmt.Println("method = ", method)
	// 打印req和reply
	fmt.Println("req = ", req.(*hello.HelloRequest).Req)
	fmt.Println("resp = ", reply.(*hello.HelloReply).Resp)
	// 打印连接目标
	fmt.Println("target = ", cc.Target())

	return err
}

func main() {
	// 创建到server的连接
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Fatal(err.Error())
	}

	defer conn.Close()
	// 初始化客户端
	client := hello.NewHelloClient(conn)
	// 使用hello方法发送消息
	resp, err := client.SayHello(context.Background(), &hello.HelloRequest{Req: "this is a request message from client"})
	if err != nil {
		log.Fatal(err.Error())
	}
	// 打印server的答复
	fmt.Println(resp.LastUpdate)
	fmt.Println(resp.Resp)
}
