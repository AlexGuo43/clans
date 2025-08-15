package main

import (
	"context"
	"log"
	"net/http"

	"github.com/AlexGuo43/clans/comment-service/config"
	"github.com/AlexGuo43/clans/comment-service/internal/handlers"
	"github.com/AlexGuo43/clans/comment-service/internal/repository"
	"github.com/AlexGuo43/clans/comment-service/internal/services"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	db := repository.ConnectDB(cfg)
	defer db.Close(context.Background())

	commentRepo := repository.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo)
	commentHandler := handlers.NewCommentHandler(commentService)

	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	api := r.PathPrefix("/api/comments").Subrouter()
	
	api.HandleFunc("", commentHandler.CreateComment).Methods("POST")
	api.HandleFunc("/{id:[0-9]+}", commentHandler.GetComment).Methods("GET")
	api.HandleFunc("/{id:[0-9]+}", commentHandler.UpdateComment).Methods("PUT")
	api.HandleFunc("/{id:[0-9]+}", commentHandler.DeleteComment).Methods("DELETE")
	api.HandleFunc("/{id:[0-9]+}/vote", commentHandler.VoteComment).Methods("POST")
	api.HandleFunc("/{id:[0-9]+}/replies", commentHandler.GetReplies).Methods("GET")
	api.HandleFunc("/post/{post_id:[0-9]+}", commentHandler.GetCommentsByPost).Methods("GET")

	log.Println("Comment Service running on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", r))
}