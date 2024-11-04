package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kannan112/gateway-structure/pkg/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	server  *http.Server
	router  *mux.Router
	logger  *zap.Logger
	options *Options
}

func NewHTTPServer(opts *Options, logger *zap.Logger) *HTTPServer {
	router := mux.NewRouter()

	server := &HTTPServer{
		router:  router,
		logger:  logger,
		options: opts,
		server: &http.Server{
			Addr:         opts.HTTPPort,
			Handler:      router,
			ReadTimeout:  opts.ReadTimeout,
			WriteTimeout: opts.WriteTimeout,
		},
	}

	server.setupRoutes()
	server.setupMiddleware()

	return server
}

func (s *HTTPServer) setupMiddleware() {
	// Add global middleware
	s.router.Use(middleware.Logger(s.logger))
	s.router.Use(middleware.Recovery())
	s.router.Use(middleware.RateLimit())
}

func (s *HTTPServer) setupRoutes() {
	//Health check
	//s.router.HandleFunc("/health", handelers.HealthCheck).Methods("GET")

	// API routes
	//api := s.router.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	//auth := api.PathPrefix("/auth").Subrouter()
	// auth.HandleFunc("/login", handlers.Login).Methods("POST")
	// auth.HandleFunc("/register", handlers.Register).Methods("POST")

	// User routes
	// 	users := api.PathPrefix("/users").Subrouter()
	// 	users.Use(middleware.Authenticate) // Protect all user routes
	// 	users.HandleFunc("", handlers.GetUsers).Methods("GET")
	// 	users.HandleFunc("/{id}", handlers.GetUser).Methods("GET")
	// 	users.HandleFunc("", handlers.CreateUser).Methods("POST")
	// 	users.HandleFunc("/{id}", handlers.UpdateUser).Methods("PUT")
}

func (s *HTTPServer) Start() error {
	s.logger.Info("Starting HTTP server", zap.String("port", s.options.HTTPPort))
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")
	return s.server.Shutdown(ctx)
}
