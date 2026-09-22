package database

import (
	"log"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/app/model"
)

func Migrate() {

	err := DB.AutoMigrate(

		&model.User{},

		&model.ElectricCarLoad{},
	)

	if err != nil {

		log.Fatal(err)

	}

	log.Println("Migration completed")

}
