package db

import (
	"musicSharingAPp/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func StartDB() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	dbURL := os.Getenv("DB_URL")
	dsn := dbURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Genre{})
	DB.AutoMigrate(&models.Playlist{})
	DB.AutoMigrate(&models.Song{})
	DB.AutoMigrate(&models.Post{})
	return nil
}
