package main

import (
	"context"
	"fmt"
	"net"

	"gitlab.ozon.ru/agarkov/route256/lectures/4-class/gRPC/interceptor/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
)

type SimpleInterceptorServer struct {
	api.UnimplementedSimpleInterceptorServer
}

func (u *SimpleInterceptorServer) UserCheck(ctx context.Context, req *api.UserRequest) (*api.UserResponse, error) {
	return &api.UserResponse{Valid: true}, nil
}

const addr = "50051"

func main() {
	listen, err := net.Listen("tcp", ":"+addr)
	if err != nil {
		return
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(EvenInterceptor),
		grpc.InTapHandle(DummyRateLimit))
	api.RegisterSimpleInterceptorServer(s, &SimpleInterceptorServer{})
	fmt.Printf("grpc server listen port:%s\n", addr)
	err = s.Serve(listen)
	if err != nil {
		panic(err)
	}
}

func EvenInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "UnaryEcho: failed to get metadata")
	}
	if v, ok := md["access-token"]; ok {
		fmt.Printf("[server-interceptor] access-token from metadata:%v\n", v)
	}

	fmt.Printf("[server-interceptor] received message info, method: %v, request:%v",
		info.FullMethod, req)
	ur := req.(*api.UserRequest)
	if ur.Id%2 != 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("[server-interceptor] only even ID allowed, reveid ID: %d", ur.Id))
	}

	m, err := handler(ctx, req)
	return m, err
}

func DummyRateLimit(ctx context.Context, info *tap.Info) (context.Context, error) {
	fmt.Println("[server-interceptor] InTapHandle - DummyRateLimit")
	return ctx, nil
}
