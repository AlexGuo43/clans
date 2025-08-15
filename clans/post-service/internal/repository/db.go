package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/AlexGuo43/clans/post-service/config"
	"github.com/jackc/pgx/v5"
)

func ConnectDB(cfg *config.Config) *pgx.Conn {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	return conn
}