package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AccessTokenMeli string
	ApiAuthToken    string
	URL             string
}

func GetEnv() (Env, error) {
	enviroment := os.Getenv("MARRIAGE_ENV")
	if enviroment == "" || enviroment == "local" {
		enviroment = "local"
	}

	fileName := ".env." + enviroment
	godotenv.Load(fileName)

	env := Env{
		AccessTokenMeli: os.Getenv("ACCESS_TOKEN"),
		ApiAuthToken:    os.Getenv("API_AUTH_TOKEN"),
		URL:             os.Getenv("URL"),
	}

	if env.AccessTokenMeli == "" {
		env.AccessTokenMeli = "DEFAULT_TOKEN"
	}
	if env.ApiAuthToken == "" {
		env.ApiAuthToken = "DEFAULT_KEY"
	}
	if env.URL == "" {
		env.URL = "0.0.0.0:8080"
	}

	log.Printf("Loaded '%s' envfile", fileName)
	return env, nil
}
