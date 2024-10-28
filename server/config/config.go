package config

import (
	"errors"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ClientURL string
	DBURL     string
	Port      string

	Email    string
	EmailPassword string
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

  email, ok := os.LookupEnv("EMAIL")
  if !ok {
    return Config{}, errors.New("error: EMAIL is not set!")
  }

  pass, ok := os.LookupEnv("EMAIL_PASSWORD")
  if !ok {
    return Config{}, errors.New("error: EMAIL_PASSWORD is not set!")
  }

  return Config{
    ClientURL: clientURL,
    DBURL: dbURL,
    Port: port,
    Email: email,
    EmailPassword: pass,
  }, nil
}
