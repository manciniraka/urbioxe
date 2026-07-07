package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string

	MailjetAPIKey    string
	MailjetSecretKey string

	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string

	OpenWeatherAPIKey string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system environment")
	}

	return &Config{
		AppPort:    os.Getenv("APP_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		MailjetAPIKey:    os.Getenv("MAILJET_API_KEY"),
		MailjetSecretKey: os.Getenv("MAILJET_SECRET_KEY"),

		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),

		OpenWeatherAPIKey: os.Getenv("OPENWEATHER_API_KEY"),
	}
}