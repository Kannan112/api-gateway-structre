package config

type Config struct {
	JWTSecret      string `mapstructure:"JWT_SRC"`
	HTTPPort       string `mapstructure:"HTTP_PORT"`
	GRPCPort       string `mapstructure:"GRPC_PORT"`
	AuthServiceURL string `mapstructure:"AUTH_SERVICE_URL"`
	UserServiceURL string `mapstructure:"USER_SERVICE_URL"`
}

var envs = []string{"JWT_SRC", "HTTP_PORT", "GRPC_PORT", "AUTH_SERVICE_URL", "USER_SERVICE_URL"}
