// @Author: Perry
// @Date  : 2020/8/28
// @Desc  : 

package main

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	hello "grpc-example/stream-example/stream-proto-out/stream-proto"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

const serverAddr = "localhost:8080"

type SHelloServer struct{}

// 客户端连续流发送，服务端持续接收然后返回单一内容
func (s *SHelloServer) CStreamHello(serv hello.StreamHello_CStreamHelloServer) error {
	var (
		err  error
		req  *hello.HelloRequest
		resp = new(hello.HelloReply)
	)
	resp.Resp = "server response"
	for {
		req, err = serv.Recv() //从客户端接收请求数据
		time.Sleep(1 * time.Second)
		if req != nil {
			fmt.Println("message received = ", req.Req)
			fmt.Println("last updated = ", time.Unix(req.LastUpdate.Seconds, 0))
		}
		if err != nil {
			fmt.Println("error = ", err.Error())
			break
		}
	}
	err = serv.SendAndClose(resp)
	return err
}

// 客户端发送单一请求，服务端持续返回数据
func (s *SHelloServer) SStreamHello(req *hello.HelloRequest, serv hello.StreamHello_SStreamHelloServer) error {
	var i = 0
	for {
		i++
		time.Sleep(1 * time.Second)
		err := serv.SendMsg(&hello.HelloReply{
			Resp:       "server reply,i=" + strconv.Itoa(i),
			LastUpdate: &timestamp.Timestamp{Seconds: time.Now().Unix()},
		})
		if err != nil {
			return err
		}
		if i == 10 {
			return nil
		}
	}
}

// 客户端持续流发送，服务端持续流返回
func (s *SHelloServer) DStreamHello(serv hello.StreamHello_DStreamHelloServer) error {
	var i int
	for {
		i++
		in, err := serv.Recv()
		if err == io.EOF {
			fmt.Println("read done")
			return nil
		}
		if err != nil {
			fmt.Println("ERR", err)
			return err
		}
		fmt.Println("userinfo ", in)
		err = serv.Send(&hello.HelloReply{Resp: strconv.Itoa(i)})
		if err != nil {
			fmt.Println("ERR", err)
		}
	}
}

func main() {
	conn, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	server := grpc.NewServer()
	hello.RegisterStreamHelloServer(server, &SHelloServer{})
	reflection.Register(server)
	log.Print("server started")
	if err = server.Serve(conn); err != nil {
		log.Fatal(err.Error())
	}
}
