package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "gitlab.ozon.ru/agarkov/route256/lectures/4-class/gRPC/gateway/gen/proto/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	httpPort = "8081"
	grpcPort = "8080"
)

type server struct {
	api.UnimplementedGreeterServer
}

func (server) SendGet(ctx context.Context, in *api.TemplateRequest) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{Message: "Received GET method " + in.Name}, nil
}

func (server) SendPost(ctx context.Context, in *api.TemplateRequest) (*api.TemplateResponse, error) {
	return &api.TemplateResponse{Message: "Received POST method " + in.Name}, nil
}
func main() {
	grpcListener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen grpc: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterGreeterServer(s, &server{})

	ctx := context.Background()
	rmux := runtime.NewServeMux()
	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	{
		err := api.RegisterGreeterHandlerServer(ctx, rmux, server{})
		if err != nil {
			log.Fatal(err)
		}
	}

	httpListener, err := net.Listen("tcp", ":"+httpPort)
	if err != nil {
		log.Fatalf("failed to listen http: %v", err)
	}
	var runHttp = func() {
		// Register reflection service on gRPC server.
		reflection.Register(s)
		if err := s.Serve(grpcListener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
	go runHttp()

	log.Printf("Serving http address %s", grpcPort)
	err = http.Serve(httpListener, mux)
	if err != nil {
		log.Fatal(err)
	}

}
