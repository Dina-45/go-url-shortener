package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Настройки подключения
	dsn := "host=localhost user=postgres password=123456 dbname=mydb port=5432 sslmode=disable"
	if os.Getenv("DB_DSN") != "" {
		dsn = os.Getenv("DB_DSN")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Connected to database!")
	DB = db
}
