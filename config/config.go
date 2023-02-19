package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Server struct {
	Port string
}

type Config struct {
	Server Server
}

func All() Config {
	loadEnv()
	return Config{
		Server: Server{Port: os.Getenv("PORT")},
	}
}

func loadEnv() {
	err := godotenv.Load("config/.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}
