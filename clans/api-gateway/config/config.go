package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ServiceConfig struct {
	Name string
	URL  string
}

type Config struct {
	Port         string
	JWTSecret    string
	UserService  ServiceConfig
	PostService  ServiceConfig
	ClanService  ServiceConfig
	Services     []ServiceConfig
}

func LoadConfig() *Config {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("Warning: No .env file found, using default values")
	}

	userService := ServiceConfig{
		Name: "user-service",
		URL:  getEnv("USER_SERVICE_URL", "http://user-service:8080"),
	}

	postService := ServiceConfig{
		Name: "post-service", 
		URL:  getEnv("POST_SERVICE_URL", "http://post-service:8081"),
	}

	clanService := ServiceConfig{
		Name: "clan-service",
		URL:  getEnv("CLAN_SERVICE_URL", "http://clan-service:8082"),
	}

	return &Config{
		Port:        getEnv("PORT", "8000"),
		JWTSecret:   getEnv("JWT_SECRET", "mysecretkey"),
		UserService: userService,
		PostService: postService,
		ClanService: clanService,
		Services: []ServiceConfig{
			userService,
			postService,
			clanService,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}