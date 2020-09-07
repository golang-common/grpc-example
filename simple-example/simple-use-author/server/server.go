// @Author: Perry
// @Date  : 2020/8/26
// @Desc  : 带拦截器的简单服务器

package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"grpc-example/simple-example/simple-proto-out/hello"
	"grpc-example/simple-example/simple-use-author/data"
	"log"
	"net"
	"strings"
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
func printInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
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
}

func tokenCheckInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		var errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
		var errInvalidToken = status.Errorf(codes.Unauthenticated, "invalid token")

		// 取token
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMissingMetadata
		}

		// 验证token
		if len(md["authorization"]) < 1 {
			return nil, errInvalidToken
		}
		token := strings.TrimPrefix(md["authorization"][0], "Bearer ")
		if token != "dpy test" {
			return nil, errInvalidToken
		}
		m, err := handler(ctx, req)
		if err != nil {
			fmt.Println(err.Error())
		}
		return m, err
	}
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
	server := grpc.NewServer(makeServerOptions()...)

	hello.RegisterHelloServer(server, &HelloServer{})
	reflection.Register(server)
	fmt.Println("server started")
	if err = server.Serve(conn); err != nil {
		log.Fatal(err.Error())
	}
}

func makeServerOptions() []grpc.ServerOption {
	var (
		srvOption []grpc.ServerOption
	)
	// unary options
	unaryOptions := grpc.ChainUnaryInterceptor(printInterceptor(), tokenCheckInterceptor())
	srvOption = append(srvOption, unaryOptions)
	// stream options
	streamOptions := grpc.ChainStreamInterceptor()
	srvOption = append(srvOption, streamOptions)
	// 证书认证
	cert, err := tls.LoadX509KeyPair(data.Path("x509/server_cert.pem"), data.Path("x509/server_key.pem"))
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}
	authOptions := grpc.Creds(credentials.NewServerTLSFromCert(&cert))
	srvOption = append(srvOption, authOptions)

	return srvOption
}
