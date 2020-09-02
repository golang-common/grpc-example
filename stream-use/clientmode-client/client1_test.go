// @Author: Perry
// @Date  : 2020/8/28
// @Desc  : 客户端流式,客户端持续发送信息，服务器单次返回信息

package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	hello "grpc-example/stream-proto-out/stream-proto"
	"strconv"
	"testing"
	"time"
)

func TestClient1(t *testing.T) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := hello.NewStreamHelloClient(conn)
	cs, err := client.CStreamHello(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		var tm = time.Now().Unix()
		var req = new(hello.HelloRequest)
		req.Req = "client send" + strconv.Itoa(i)
		req.LastUpdate = &timestamp.Timestamp{Seconds: tm}
		err = cs.Send(req)
		if err != nil {
			t.Fatal(err)
		}
	}
	var reply = new(hello.HelloReply)
	err = cs.CloseSend()
	if err != nil {
		t.Fatal(err)
	}
	err = cs.RecvMsg(reply)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply.Resp)
	t.Log(reply.LastUpdate)
}
