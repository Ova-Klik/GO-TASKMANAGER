package db

import (
	"log"

	"TaskManager/config"
	"TaskManager/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg config.Config) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	if err = DB.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Database connected and migrated successfully.")
}
