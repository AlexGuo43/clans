package main

import (
	"context"
	"log"
	"net/http"

	"github.com/AlexGuo43/clans/post-service/config"
	"github.com/AlexGuo43/clans/post-service/internal/handlers"
	"github.com/AlexGuo43/clans/post-service/internal/repository"
	"github.com/AlexGuo43/clans/post-service/internal/services"
	"github.com/gorilla/mux"
)


func main() {
	cfg := config.LoadConfig()
	db := repository.ConnectDB(cfg)
	defer db.Close(context.Background())

	postRepo := repository.NewPostRepository(db)
	postService := services.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	api := r.PathPrefix("/api/posts").Subrouter()
	
	api.HandleFunc("", postHandler.GetPosts).Methods("GET")
	api.HandleFunc("/{id:[0-9]+}", postHandler.GetPost).Methods("GET")
	api.HandleFunc("/clan/{clan_id:[0-9]+}", postHandler.GetPostsByClan).Methods("GET")

	api.HandleFunc("", postHandler.CreatePost).Methods("POST")
	api.HandleFunc("/{id:[0-9]+}", postHandler.UpdatePost).Methods("PUT")
	api.HandleFunc("/{id:[0-9]+}", postHandler.DeletePost).Methods("DELETE")
	api.HandleFunc("/{id:[0-9]+}/vote", postHandler.VotePost).Methods("POST")

	log.Println("Post Service running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", r))
}