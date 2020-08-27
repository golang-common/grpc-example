// @Author: Perry
// @Date  : 2020/8/26
// @Desc  : 

package main

import (
	"context"
	"fmt"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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
	server := grpc.NewServer()
	hello.RegisterHelloServer(server, &HelloServer{})
	reflection.Register(server)
	fmt.Println("server started")
	if err = server.Serve(conn); err != nil {
		log.Fatal(err.Error())
	}
}
