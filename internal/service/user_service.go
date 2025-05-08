package service

import (
	"context"
	"grpc-gateway-demo/proto/gen/go/pcgw"
	"log"
)

type UserService struct {
	pcgw.UnimplementedUserServiceServer
}

func (u *UserService) AddUser(ctx context.Context, request *pcgw.AddUserRequest) (*pcgw.AddUserResponse, error) {
	log.Printf("Received request for adding a new user: %+v", request)
	return &pcgw.AddUserResponse{
		Id: "user-id",
	}, nil
}
