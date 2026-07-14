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

	MailjetBaseURL     string
	MailjetAPIKey      string
	MailjetSecretKey   string
	MailjetSenderEmail string
	MailjetSenderName  string

	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string

	OpenWeatherBaseURL string
	OpenWeatherAPIKey  string
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

		MailjetBaseURL:     os.Getenv("MAILJET_BASE_URL"),
		MailjetAPIKey:      os.Getenv("MAILJET_API_KEY"),
		MailjetSecretKey:   os.Getenv("MAILJET_SECRET_KEY"),
		MailjetSenderEmail: os.Getenv("MAILJET_SENDER_EMAIL"),
		MailjetSenderName:  os.Getenv("MAILJET_SENDER_NAME"),

		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),

		OpenWeatherBaseURL: os.Getenv("OPENWEATHER_BASE_URL"),
		OpenWeatherAPIKey:  os.Getenv("OPENWEATHER_API_KEY"),
	}
}
