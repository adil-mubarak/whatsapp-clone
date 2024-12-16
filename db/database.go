package db

import (
	"fmt"
	"log"
	"os"
	"whatsapp/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDB() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB,err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil{
		log.Fatalf("Failed to connect to database: %v",err)
	}

	log.Println("Database connected successfully")

	if err := DB.AutoMigrate(models.User{}); err != nil{
		log.Println("failed to migreate database: ",err)
	}
}
