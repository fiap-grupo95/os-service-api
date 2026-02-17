package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDBConfig struct {
	URI      string
	Database string
	Client   *mongo.Client
}

func NewMongoDBFromEnv() (*MongoDBConfig, error) {
	uri := os.Getenv("MONGODB_URI")
	database := os.Getenv("MONGODB_DATABASE")

	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable is not set")
	}

	if database == "" {
		database = "os-service-db"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(50).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second).
		SetReadPreference(readpref.SecondaryPreferred())

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logs.Logger().Info().Str("database", database).Msg("Successfully connected to MongoDB")

	return &MongoDBConfig{
		URI:      uri,
		Database: database,
		Client:   client,
	}, nil
}

func (c *MongoDBConfig) GetDatabase() *mongo.Database {
	return c.Client.Database(c.Database)
}

func (c *MongoDBConfig) GetCollection(collectionName string) *mongo.Collection {
	return c.GetDatabase().Collection(collectionName)
}

func (c *MongoDBConfig) Disconnect(ctx context.Context) error {
	if c.Client != nil {
		return c.Client.Disconnect(ctx)
	}
	return nil
}
