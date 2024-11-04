package server

import (
	"time"

	"github.com/kannan112/gateway-structure/pkg/config"
	"github.com/kannan112/gateway-structure/pkg/service"
)

type Options struct {
	HTTPPort        string
	GRPCPort        string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	AuthService     service.AuthServiceConfig
	UserService     service.UserServiceConfig
}

func DefaultOptions(conf *config.Config) *Options {
	return &Options{
		HTTPPort:        ":8080",
		GRPCPort:        ":9090",
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		AuthService: service.AuthServiceConfig{
			Address: "localhost:50052",
			Timeout: 10 * time.Second,
		},
	}
}
