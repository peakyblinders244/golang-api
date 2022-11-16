package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Env struct {
	CLIENT_URL string
	DB_USER    string
	DB_PASS    string
	DB_HOST    string
	DB_NAME    string
	DB_PORT    string
	JWT_SECRET string
	EMAIL_FROM string
	SMTP_HOST  string
	SMTP_PORT  string
	SMTP_USER  string
	SMTP_PASS  string
	PORT       string
}

func LoadEnv() Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return Env{
		CLIENT_URL: os.Getenv("CLIENT_URL"),
		DB_USER:    os.Getenv("DB_USER"),
		DB_PASS:    os.Getenv("DB_PASS"),
		DB_HOST:    os.Getenv("DB_HOST"),
		DB_NAME:    os.Getenv("DB_NAME"),
		DB_PORT:    os.Getenv("DB_PORT"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
		EMAIL_FROM: os.Getenv("EMAIL_FROM"),
		SMTP_HOST:  os.Getenv("SMTP_HOST"),
		SMTP_PORT:  os.Getenv("SMTP_PORT"),
		SMTP_USER:  os.Getenv("SMTP_USER"),
		SMTP_PASS:  os.Getenv("SMTP_PASS"),
		PORT:       os.Getenv("PORT"),
	}
}
