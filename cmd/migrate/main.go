package main

import (
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/database"
)

func main() {
	db := database.ConnectDatabase()
	database.Migrate(db)
	database.Seed(db)
}
