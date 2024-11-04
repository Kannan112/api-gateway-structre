package handlers

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/fluxor-api-gateway/pkg/service"

	userpb "github.com/fluxor-api-gateway/pkg/proto/user"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userClient service.UserService
	logger     *log.Logger
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userClient service.UserService, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userClient: userClient,
		logger:     logger,
	}
}

// CreateUser handles user creation requests
func (h *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	if err := validateCreateUserRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response, err := h.userClient.CreateUser(ctx, req)
	if err != nil {
		h.logger.Printf("User creation failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return response, nil
}

// GetUser retrieves user information
func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	if err := validateGetUserRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response, err := h.userClient.GetUser(ctx, req)
	if err != nil {
		h.logger.Printf("Get user failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return response, nil
}

// UpdateUser handles user update requests
func (h *UserHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	if err := validateUpdateUserRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response, err := h.userClient.UpdateUser(ctx, req)
	if err != nil {
		h.logger.Printf("User update failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return response, nil
}

// DeleteUser handles user deletion requests
func (h *UserHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	if err := validateDeleteUserRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response, err := h.userClient.DeleteUser(ctx, req)
	if err != nil {
		h.logger.Printf("User deletion failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	return response, nil
}

// ListUsers retrieves a list of users with pagination
func (h *UserHandler) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	if err := validateListUsersRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response, err := h.userClient.ListUsers(ctx, req)
	if err != nil {
		h.logger.Printf("List users failed: %v", err)
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	return response, nil
}

// Helper functions for request validation
func validateCreateUserRequest(req *userpb.CreateUserRequest) error {
	if req.User == nil {
		return status.Error(codes.InvalidArgument, "user information is required")
	}
	if req.User.Email == "" || req.User.Username == "" {
		return status.Error(codes.InvalidArgument, "email and username are required")
	}
	return nil
}

func validateGetUserRequest(req *userpb.GetUserRequest) error {
	if req.UserId == "" {
		return status.Error(codes.InvalidArgument, "user ID is required")
	}
	return nil
}

func validateUpdateUserRequest(req *userpb.UpdateUserRequest) error {
	if req.User == nil {
		return status.Error(codes.InvalidArgument, "user information is required")
	}
	if req.User.Id == "" {
		return status.Error(codes.InvalidArgument, "user ID is required")
	}
	return nil
}

func validateDeleteUserRequest(req *userpb.DeleteUserRequest) error {
	if req.UserId == "" {
		return status.Error(codes.InvalidArgument, "user ID is required")
	}
	return nil
}

func validateListUsersRequest(req *userpb.ListUsersRequest) error {
	if req.PageSize < 0 {
		return status.Error(codes.InvalidArgument, "page size must be non-negative")
	}
	return nil
}
