package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/AlexGuo43/clans/clan-service/internal/config"
	"github.com/AlexGuo43/clans/clan-service/internal/handlers"
	"github.com/AlexGuo43/clans/clan-service/internal/repository"
	"github.com/AlexGuo43/clans/clan-service/internal/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.LoadConfig()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database successfully")

	clanRepo := repository.NewClanRepository(db)
	clanService := services.NewClanService(clanRepo)
	clanHandler := handlers.NewClanHandler(clanService)

	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/clans", clanHandler.CreateClan).Methods("POST")
	api.HandleFunc("/clans", clanHandler.GetClans).Methods("GET")
	api.HandleFunc("/clans/{id:[0-9]+}", clanHandler.GetClan).Methods("GET")
	api.HandleFunc("/clans/name/{name}", clanHandler.GetClanByName).Methods("GET")
	api.HandleFunc("/clans/{id:[0-9]+}", clanHandler.UpdateClan).Methods("PUT")
	api.HandleFunc("/clans/{id:[0-9]+}", clanHandler.DeleteClan).Methods("DELETE")
	
	api.HandleFunc("/clans/{id:[0-9]+}/join", clanHandler.JoinClan).Methods("POST")
	api.HandleFunc("/clans/{id:[0-9]+}/leave", clanHandler.LeaveClan).Methods("POST")
	api.HandleFunc("/clans/{id:[0-9]+}/members", clanHandler.GetMembers).Methods("GET")
	api.HandleFunc("/clans/{clanId:[0-9]+}/members/{userId:[0-9]+}/role", clanHandler.UpdateMemberRole).Methods("PUT")
	api.HandleFunc("/clans/{id:[0-9]+}/membership", clanHandler.GetMembership).Methods("GET")
	
	api.HandleFunc("/users/clans", clanHandler.GetUserClans).Methods("GET")

	log.Printf("Clan service starting on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, r))
}