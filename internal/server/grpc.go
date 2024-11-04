package server

import (
	"fmt"
	"net"

	"github.com/kannan112/gateway-structure/pkg/middleware"
	"github.com/kannan112/gateway-structure/pkg/proto/auth"
	"github.com/kannan112/gateway-structure/pkg/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server  *grpc.Server
	logger  *zap.Logger
	options *Options
}

func NewGRPCServer(opts *Options, logger *zap.Logger) (*GRPCServer, error) {
	// Create gRPC server with interceptors
	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.GRPCLogger(logger)),
		grpc.ChainUnaryInterceptor(
			middleware.GRPCRecovery(),
			middleware.GRPCAuth(),
		),
	)

	// Initialize services with error handling
	authService, err := service.NewAuthService(opts.AuthService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth service: %v", err)
	}

	// userService, err := service.NewUserService(opts.UserService)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to initialize user service: %v", err)
	// }

	// Register services
	auth.RegisterAuthServiceServer(server, authService)
	//user.RegisterUserServiceServer(server, userService)

	// Enable reflection for grpcurl
	reflection.Register(server)

	return &GRPCServer{
		server:  server,
		logger:  logger,
		options: opts,
	}, nil
}

func (s *GRPCServer) Start() error {
	listener, err := net.Listen("tcp", s.options.GRPCPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.logger.Info("Starting gRPC server", zap.String("port", s.options.GRPCPort))
	return s.server.Serve(listener)
}

func (s *GRPCServer) Stop() {
	s.logger.Info("Stopping gRPC server")
	s.server.GracefulStop()
}
