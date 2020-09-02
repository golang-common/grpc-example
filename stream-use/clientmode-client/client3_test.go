// @Author: Perry
// @Date  : 2020/8/28
// @Desc  : 双向流模式，客户端和服务器

package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	hello "grpc-example/stream-proto-out/stream-proto"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestClient3(t *testing.T) {
	wg := sync.WaitGroup{}
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := hello.NewStreamHelloClient(conn)
	cs, err := client.DStreamHello(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			var tm = time.Now().Unix()
			var req = new(hello.HelloRequest)
			req.Req = "client send" + strconv.Itoa(i)
			req.LastUpdate = &timestamp.Timestamp{Seconds: tm}
			err = cs.SendMsg(req)
			if err != nil {
				t.Log(err)
			}
		}
		cs.CloseSend()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			rpl, err := cs.Recv()
			if err != nil {
				t.Log(err)
				break
			}
			fmt.Println("response from server = ", rpl.Resp)
		}
	}()
	wg.Wait()
}
