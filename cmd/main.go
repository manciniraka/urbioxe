package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/database"
	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/router"
)

func main() {
	cfg := config.Load()

	db := database.ConnectPostgres(cfg)

	_ = db

	e := echo.New()

	router.InitRouter(e)

	log.Printf("Starting ubioxe API...")
	log.Printf("ubioxe server running on :%s\n", cfg.AppPort)

	if err := e.Start(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		log.Fatal(err)
	}
}
