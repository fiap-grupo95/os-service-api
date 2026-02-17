package main

import (
	"context"
	"log"

	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/database"
)

func main() {
	ctx := context.Background()
	db := database.ConnectDatabase()

	if err := database.SeedMongo(ctx, db); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	log.Println("seed completed")
}
