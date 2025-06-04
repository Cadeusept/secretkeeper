package grpc

import (
	"context"

	"github.com/Cadeusept/secretkeeper/internal/usecases/secretkeeper"
	pb "github.com/Cadeusept/secretkeeper/proto/secretkeeper"
)

type SecretKeeperServer struct {
	pb.UnimplementedSecretKeeperServer
	usc *secretkeeper.Usc
}

func NewSecretKeeperServer(usc *secretkeeper.Usc) *SecretKeeperServer {
	return &SecretKeeperServer{usc: usc}
}

func (s *SecretKeeperServer) AddKey(ctx context.Context, req *pb.AddKeyRequest) (*pb.AddKeyResponse, error) {
	err := s.usc.AddKey(ctx, req.UserId, req.ServiceId, req.ApiKey)
	if err != nil {
		return nil, err
	}
	return &pb.AddKeyResponse{Success: true}, nil
}

func (s *SecretKeeperServer) GetKey(ctx context.Context, req *pb.GetKeyRequest) (*pb.GetKeyResponse, error) {
	apiKey, err := s.usc.GetKey(ctx, req.UserId, req.ServiceId)
	if err != nil {
		return nil, err
	}
	return &pb.GetKeyResponse{ApiKey: apiKey}, nil
}

func (s *SecretKeeperServer) UpdateKey(ctx context.Context, req *pb.UpdateKeyRequest) (*pb.UpdateKeyResponse, error) {
	err := s.usc.UpdateKey(ctx, req.UserId, req.ServiceId, req.NewApiKey)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateKeyResponse{Success: true}, nil
}
