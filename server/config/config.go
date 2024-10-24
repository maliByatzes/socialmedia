package config

import (
	"errors"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
  ClientURL string
  DBURL string
  Port string
}

func NewConfig() (Config, error) {
  clientURL, ok := os.LookupEnv("CLIENT_URL")
  if !ok {
    return Config{}, errors.New("error: CLIENT_URL is not set!")
  }

  dbURL, ok := os.LookupEnv("DB_URL")
  if !ok {
    return Config{}, errors.New("error: DB_URL is not set!")
  }

  port, ok := os.LookupEnv("PORT")
  if !ok {
    port = "8080"
  }

  return Config{ClientURL: clientURL, DBURL: dbURL, Port: port}, nil
}
