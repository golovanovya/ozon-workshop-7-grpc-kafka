package main

import (
	"context"
	"fmt"
	"log"

	"gitlab.ozon.ru/agarkov/route256/lectures/4-class/gRPC/interceptor/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const addr = "50051"

func main() {
	dial, err := grpc.Dial(":"+addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(SimpleUnaryInterceptor),
		grpc.WithPerRPCCredentials(&token{Secret: "9230750293"}))
	if err != nil {
		panic(err)
	}
	defer dial.Close()

	client := api.NewSimpleInterceptorClient(dial)

	check, err := client.UserCheck(context.Background(), &api.UserRequest{
		Id:    2,
		Login: "Denis",
	})
	if err != nil {
		errStatus, _ := status.FromError(err)
		fmt.Printf("[client-handler]%v\n", errStatus.Message())
		fmt.Printf("[client-handler]%v\n", errStatus.Code())

		if codes.InvalidArgument == errStatus.Code() {
			log.Fatal()
		}
	}
	fmt.Printf("[client-handler] %s", check.String())
}

// SimpleUnaryInterceptor is a gRPC UnaryServerInterceptor that will greeting.
func SimpleUnaryInterceptor(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Printf("[client-interceptor] Hello from unary client interceptor, request: %v\n", req)
	return invoker(ctx, method, req, reply, cc, opts...)
}

type token struct {
	Secret string
}

func (t *token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"access-token": t.Secret,
	}, nil
}

func (t *token) RequireTransportSecurity() bool {
	return false
}
