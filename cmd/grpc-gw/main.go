package main

import (
	"context"
	"grpc-gateway-demo/internal/service"
	"grpc-gateway-demo/proto/gen/go/gateway"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	go runGRPC()
	runHTTP()
}

func runGRPC() {
	listener, err := net.Listen("tcp", ":40000")
	if err != nil {
		log.Fatalf("Failed to bind: %s", err)
	}

	server := grpc.NewServer()
	gateway.RegisterUserServiceServer(server, &service.UserService{})
	reflection.Register(server)

	log.Printf("Starting gRPC server on port 40000...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %s", err)
	}
}

func runHTTP() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := gateway.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:40000", opts)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("REST gateway listening on :8888")
	if err := http.ListenAndServe(":8888", mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
