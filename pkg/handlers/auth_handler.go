package handlers

import (
	"log"

	"github.com/kannan112/gateway-structure/pkg/service"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	authClient service.AuthService
	logger     *log.Logger
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authClient service.AuthService, logger *log.Logger) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
		logger:     logger,
	}
}
