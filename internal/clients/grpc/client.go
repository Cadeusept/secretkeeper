package grpc

import (
	"context"
	"time"

	pb "github.com/Cadeusept/secretkeeper/proto/secretkeeper"

	"google.golang.org/grpc"
)

type SecretClient struct {
	conn   *grpc.ClientConn
	client pb.SecretKeeperClient
}

// NewSecretClient initializes a new SecretClient
func NewSecretClient(address string) (*SecretClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	client := pb.NewSecretKeeperClient(conn)
	return &SecretClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection
func (c *SecretClient) Close() error {
	return c.conn.Close()
}

// AddKey adds a new API key for a specific user and service
func (c *SecretClient) AddKey(ctx context.Context, userID, serviceID, apiKey string) (bool, error) {
	resp, err := c.client.AddKey(ctx, &pb.AddKeyRequest{
		UserId:    userID,
		ServiceId: serviceID,
		ApiKey:    apiKey,
	})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

// GetKey retrieves the API key for a specific user and service
func (c *SecretClient) GetKey(ctx context.Context, userID, serviceID string) (string, error) {
	resp, err := c.client.GetKey(ctx, &pb.GetKeyRequest{
		UserId:    userID,
		ServiceId: serviceID,
	})
	if err != nil {
		return "", err
	}
	return resp.ApiKey, nil
}

// UpdateKey updates the API key for a specific user and service
func (c *SecretClient) UpdateKey(ctx context.Context, userID, serviceID, newKey string) (bool, error) {
	resp, err := c.client.UpdateKey(ctx, &pb.UpdateKeyRequest{
		UserId:    userID,
		ServiceId: serviceID,
		NewApiKey: newKey,
	})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}
