package main

import (
	"log"

	"github.com/Korsar21/electric-vehicle-winter-range-backend/internal/api"
)

func main() {
	log.Println("Application start")

	if err := api.StartServer(); err != nil {
		log.Fatal(err)
	}
}
