package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fluxor-api-gateway/internal/server"
	"github.com/fluxor-api-gateway/pkg/config"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger'
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config %s", err)
		return
	}

	// Create server options
	opts := server.DefaultOptions(&config)

	// Initialize servers
	httpServer := server.NewHTTPServer(opts, logger)
	grpcServer, err := server.NewGRPCServer(opts, logger)

	if err != nil {
		logger.Fatal("Failed to initialize gRPC server",
			zap.Error(err),
			zap.String("auth_service", opts.AuthService.Address),
			zap.String("grpc_port", opts.GRPCPort),
		)
		os.Exit(1)
	}

	// Start servers in goroutines
	go func() {
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	go func() {
		if err := grpcServer.Start(); err != nil {
			logger.Fatal("Failed to start gRPC server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), opts.ShutdownTimeout)
	defer cancel()

	// Stop servers
	if err := httpServer.Stop(ctx); err != nil {
		logger.Error("Failed to stop HTTP server", zap.Error(err))
	}
	grpcServer.Stop()

	logger.Info("Servers stopped successfully")
}
