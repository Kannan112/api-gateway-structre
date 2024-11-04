package service

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/kannan112/gateway-structure/pkg/proto/auth"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	authpb.AuthServiceServer
	Close() error
}

// AuthServiceConfig holds configuration for the auth service client
type AuthServiceConfig struct {
	Address string
	Timeout time.Duration
}

// authServiceServer implements AuthService interface
type authServiceServer struct {
	authpb.UnimplementedAuthServiceServer
	client  authpb.AuthServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(config AuthServiceConfig) (AuthService, error) {
	if config.Address == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}

	// Set default timeout
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Connection options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithReturnConnectionError(), // This helps with more detailed error messages
	}

	// Establish gRPC connection
	conn, err := grpc.DialContext(ctx, config.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service at %s: %v", config.Address, err)
	}

	// Verify connection state
	state := conn.GetState()
	if state != connectivity.Ready {
		conn.Close()
		return nil, fmt.Errorf("connection not ready, current state: %v", state)
	}

	return &authServiceServer{
		client:  authpb.NewAuthServiceClient(conn),
		conn:    conn,
		timeout: config.Timeout,
	}, nil
}

// ensureConnection checks and ensures the connection is healthy
func (s *authServiceServer) ensureConnection(ctx context.Context) error {
	state := s.conn.GetState()
	if state == connectivity.TransientFailure || state == connectivity.Shutdown {
		return fmt.Errorf("connection is in %v state", state)
	}
	return nil
}

// Close closes the gRPC connection
func (s *authServiceServer) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
