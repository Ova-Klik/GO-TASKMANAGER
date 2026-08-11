package main

import (
	"log"

	"TaskManager/config"
	"TaskManager/db"
	"TaskManager/middleware"
	"TaskManager/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not loaded, falling back to system environment variables")
	}

	cfg := config.Load()

	db.Connect(cfg)

	r := gin.Default()
	r.Use(middleware.Logger())

	routes.SetupRoutes(r)

	log.Printf("Server starting on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
