// @Author: Perry
// @Date  : 2020/8/26
// @Desc  : 

package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"grpc-example/simple-example/simple-proto-out/hello"
	"grpc-example/simple-example/simple-use-author/data"
	"log"
)

const serverAddr = "localhost:8080"

func main() {
	// 创建到server的连接
	conn, err := grpc.Dial(serverAddr, makeDialOptions()...)
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

func makeDialOptions() []grpc.DialOption {
	var dialOption []grpc.DialOption
	// 设置证书加密
	creds, err := credentials.NewClientTLSFromFile(data.Path("x509/ca_cert.pem"), "x.test.example.com")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	dialOption = append(dialOption, grpc.WithTransportCredentials(creds))
	// 设置token
	tk := &oauth2.Token{AccessToken: "dpy test"}
	rpcRredential := oauth.NewOauthAccess(tk)
	tokenOption := grpc.WithPerRPCCredentials(rpcRredential)
	// 插入token验证
	dialOption = append(dialOption, tokenOption)
	return dialOption
}
