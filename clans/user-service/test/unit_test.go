package repository_test

import (
	"context"
	"testing"

	"github.com/AlexGuo43/clans/user-service/internal/models"
	"github.com/AlexGuo43/clans/user-service/internal/repository"
	"github.com/jackc/pgx/v5"
)

func TestCreateUser(t *testing.T) {
	// Mock or connect to a test database
	db, _ := pgx.Connect(context.Background(), "postgres://admin:adminpass@localhost:5432/clans")
	repo := repository.UserRepository{DB: db}

	user := &models.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "hashedpassword",
	}

	err := repo.CreateUser(user)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}
}
