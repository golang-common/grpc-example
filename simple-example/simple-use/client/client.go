// @Author: Perry
// @Date  : 2020/8/26
// @Desc  : 

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-example/simple-example/simple-proto-out/hello"
	"log"
)

const serverAddr = "localhost:8080"

func main() {
	// 创建到server的连接
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
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
