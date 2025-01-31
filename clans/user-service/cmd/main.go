package main

import (
	"context"
	"log"
	"net/http"

	"github.com/AlexGuo43/clans/user-service/config"
	"github.com/AlexGuo43/clans/user-service/internal/handlers"
	"github.com/AlexGuo43/clans/user-service/internal/middleware"
	"github.com/AlexGuo43/clans/user-service/internal/repository"
	"github.com/AlexGuo43/clans/user-service/internal/services"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	db := repository.ConnectDB(cfg)
	defer db.Close(context.Background())

	// Initialize services and handlers
	userRepo := &repository.UserRepository{DB: db}
	userService := &services.UserService{Repo: userRepo}
	userHandler := &handlers.UserHandler{UserService: userService}

	// Set up routes
	r := mux.NewRouter()
	r.HandleFunc("/signup", userHandler.RegisterUser).Methods("POST")
	r.HandleFunc("/login", userHandler.LoginUser).Methods("POST")

	// Protected route (requires authentication)
	protected := r.PathPrefix("/protected").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the protected dashboard!"))
	}).Methods("GET")

	// Start the server
	log.Println("User Service running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
