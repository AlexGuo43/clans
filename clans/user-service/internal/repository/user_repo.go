package repository

import (
	"context"
	"errors"

	"github.com/AlexGuo43/clans/user-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	DB *pgx.Conn
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	_, err := repo.DB.Exec(context.Background(),
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3)",
		user.Username, user.Email, user.Password)
	return err
}

func (repo *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := repo.DB.QueryRow(context.Background(),
		"SELECT id, username, email, password FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}
