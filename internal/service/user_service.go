package service

import (
	"context"
	"grpc-gateway-demo/proto/gen/go/gateway"
	"log"
)

type UserService struct {
	gateway.UnimplementedUserServiceServer
}

func (u *UserService) AddUser(ctx context.Context, request *gateway.AddUserRequest) (*gateway.AddUserResponse, error) {
	log.Printf("Received request for adding a new user: %+v", request)
	return &gateway.AddUserResponse{
		Id: "user-id",
	}, nil
}
