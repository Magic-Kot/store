package grpcapp

import (
	"context"

	"google.golang.org/grpc"

	storev1 "github.com/Magic-Kot/store-api/gen/go"
)

type AuthServer struct {
	storev1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	storev1.RegisterAuthServer(gRPC, &AuthServer{})
}

func (a AuthServer) Register(ctx context.Context, req *storev1.RegisterRequest) (*storev1.RegisterResponse, error) {
	return nil, nil
}

func (a AuthServer) Login(ctx context.Context, req *storev1.LoginRequest) (*storev1.LoginResponse, error) {
	return nil, nil
}
