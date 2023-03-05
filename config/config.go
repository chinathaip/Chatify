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
	Server       Server
	DBConnection string
}

func All() Config {
	loadEnv()
	return Config{
		Server:       Server{Port: os.Getenv("PORT")},
		DBConnection: os.Getenv("DB_CONNECTION"),
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}
