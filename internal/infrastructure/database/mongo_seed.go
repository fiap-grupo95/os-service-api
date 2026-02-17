package database

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func SeedMongo(ctx context.Context, db *MongoDBConfig) error {
	defaultPassword := os.Getenv("MONGO_SEED_DEFAULT_PASSWORD")
	if defaultPassword == "" {
		return errors.New("MONGO_SEED_DEFAULT_PASSWORD is required to seed users")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	coll := db.GetCollection("users")
	now := time.Now()

	users := []bson.M{
		{"user_id": "1", "email": "admin@xpto.com", "password": string(hashedPassword), "user_type": "admin", "created_at": now, "updated_at": now},
		{"user_id": "2", "email": "joao@xpto.com", "password": string(hashedPassword), "user_type": "customer", "created_at": now, "updated_at": now},
		{"user_id": "3", "email": "joana@xpto.com", "password": string(hashedPassword), "user_type": "customer", "created_at": now, "updated_at": now},
	}

	for _, u := range users {
		filter := bson.M{"email": u["email"], "deleted_at": bson.M{"$exists": false}}
		update := bson.M{"$setOnInsert": u}
		_, err := coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				continue
			}
			return err
		}
	}

	return nil
}
