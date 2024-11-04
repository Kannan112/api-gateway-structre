package service

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "github.com/kannan112/gateway-structure/pkg/proto/user"
)

// UserService defines the interface for user operations
type UserService interface {
	userpb.UserServiceServer
	Close() error
}

// userServiceServer implements UserService interface
type userServiceServer struct {
	userpb.UnimplementedUserServiceServer // Embed the unimplemented server
	client                                userpb.UserServiceClient
	conn                                  *grpc.ClientConn
	timeout                               time.Duration
}

// UserServiceConfig holds configuration for the user service client
type UserServiceConfig struct {
	Address string
	Timeout time.Duration
}

// NewUserService creates a new instance of UserService
func NewUserService(config UserServiceConfig) (UserService, error) {
	// Set default timeout if not provided
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Establish gRPC connection
	conn, err := grpc.Dial(
		config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}

	return &userServiceServer{
		client:  userpb.NewUserServiceClient(conn),
		conn:    conn,
		timeout: config.Timeout,
	}, nil
}

// CreateUser implements the user creation operation
func (s *userServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.client.CreateUser(ctx, req)
}

// GetUser implements the get user operation
func (s *userServiceServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.client.GetUser(ctx, req)
}

// UpdateUser implements the user update operation
func (s *userServiceServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.client.UpdateUser(ctx, req)
}

// DeleteUser implements the user deletion operation
func (s *userServiceServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.client.DeleteUser(ctx, req)
}

// ListUsers implements the list users operation
func (s *userServiceServer) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.client.ListUsers(ctx, req)
}

// Close closes the gRPC connection
func (s *userServiceServer) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
