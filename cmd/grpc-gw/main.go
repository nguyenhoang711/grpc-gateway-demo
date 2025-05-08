package main

import (
	"context"
	"fmt"
	"grpc-gateway-demo/internal/service"
	"grpc-gateway-demo/internal/websocket"
	"grpc-gateway-demo/internal/websocket/wrapper"
	"grpc-gateway-demo/proto/gen/go/pcgw"
	"log"
	"net"
	"net/http"

	ws "github.com/gorilla/websocket"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":40000")
	if err != nil {
		log.Fatalf("failed to bind: %v", err)
	}

	server := grpc.NewServer()
	pcgw.RegisterUserServiceServer(server, &service.UserService{})
	reflection.Register(server)

	connection, err := grpc.NewClient(":40000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	upgrader := ws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	_ = func(w http.ResponseWriter, r *http.Request) (websocket.Connection, error) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("ws error: %s", err)
			return nil, fmt.Errorf("failed to upgrade: %w", err)
		}

		return wrapper.New(conn), nil
	}

	restGateway := runtime.NewServeMux()
	if err := pcgw.RegisterUserServiceHandlerClient(context.Background(), restGateway, pcgw.NewUserServiceClient(connection)); err != nil {
		log.Fatal(err)
	}
	httpMux := http.NewServeMux()
	httpMux.Handle("/", restGateway)

	go func() {
		log.Printf("starting HTTP on port 8888...")
		if err := http.ListenAndServe(":8888", httpMux); err != nil {
			log.Fatalf("failed to start HTTP Rest Gateway service: %s", err)
		}
	}()

	log.Printf("starting gRPC on port 40000...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to start gRPC server: %s", err)
	}
}

func runGRPC() {
	listener, err := net.Listen("tcp", ":40000")
	if err != nil {
		log.Fatalf("Failed to bind: %s", err)
	}

	server := grpc.NewServer()
	pcgw.RegisterUserServiceServer(server, &service.UserService{})
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

	err := pcgw.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:40000", opts)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("REST gateway listening on :8888")
	if err := http.ListenAndServe(":8888", mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
