package main

import (
	"log"
	"net/http"

	"github.com/AlexGuo43/clans/api-gateway/config"
	"github.com/AlexGuo43/clans/api-gateway/internal/middleware"
	"github.com/AlexGuo43/clans/api-gateway/internal/proxy"
	"github.com/AlexGuo43/clans/api-gateway/internal/services"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	authService := services.NewAuthService(cfg.JWTSecret)
	gateway := proxy.NewGateway(cfg)

	r := mux.NewRouter()

	r.HandleFunc("/health", gateway.HealthCheck).Methods("GET")
	
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.LoggingMiddleware)
	api.Use(middleware.CorsMiddleware)
	api.Use(middleware.AuthMiddleware(authService))
	
	api.PathPrefix("/").HandlerFunc(gateway.RouteRequest)

	log.Printf("API Gateway starting on port %s...", cfg.Port)
	log.Printf("Routing to services:")
	for _, service := range cfg.Services {
		log.Printf("  - %s: %s", service.Name, service.URL)
	}

	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}