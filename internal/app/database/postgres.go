package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	dsn := fmt.Sprintf(
		"host=127.0.0.1 user=postgres password=postgres dbname=electric_car_loads port=5433 sslmode=disable TimeZone=Europe/Moscow",
	)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	DB = db

	log.Println("PostgreSQL connected")
}
