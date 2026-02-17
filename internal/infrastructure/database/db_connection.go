package database

import (
	"log"
)

func ConnectDatabase() *MongoDBConfig {
	db, err := NewMongoDBFromEnv()
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	return db
}
