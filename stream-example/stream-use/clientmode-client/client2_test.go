// @Author: Perry
// @Date  : 2020/8/28
// @Desc  : 服务器流模式，客户端单次请求，服务器持续返回信息

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	hello "grpc-example/stream-example/stream-proto-out/stream-proto"
	"testing"
	"time"
)

func TestClient2(t *testing.T) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := hello.NewStreamHelloClient(conn)
	// 发送请求信息
	cs, err := client.SStreamHello(context.Background(), &hello.HelloRequest{Req: "client req sent"})
	if err != nil {
		t.Fatal(err)
	}
	// 从服务器接收信息
	for {
		reply, err := cs.Recv()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(reply.Resp)
		fmt.Println(time.Unix(reply.LastUpdate.Seconds, 0))
	}
}
