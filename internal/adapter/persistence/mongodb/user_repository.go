package mongodb

import (
	"context"
	"time"

	"github.com/fiap-grupo95/os-service-api/internal/domain/entities"
	"github.com/fiap-grupo95/os-service-api/internal/domain/valueobject"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/database"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/usecase/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserMongoDB struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
	UserType  string             `bson:"user_type" json:"user_type"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type UserRepository struct {
	collection *mongo.Collection
}

var _ interfaces.IUserRepository = (*UserRepository)(nil)

func NewUserRepository(db *database.MongoDBConfig) *UserRepository {
	repo := &UserRepository{
		collection: db.GetCollection("users"),
	}
	_ = repo.createIndexes(context.Background())
	return repo
}

func (r *UserRepository) createIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *UserRepository) GetByID(id string) (*entities.User, error) {
	logger := logs.Logger()

	var mongoModel UserMongoDB
	filter := bson.M{"user_id": id, "deleted_at": bson.M{"$exists": false}}

	err := r.collection.FindOne(context.Background(), filter).Decode(&mongoModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logger.Error().Err(err).Str("user_id", id).Msg("Error finding user")
		return nil, err
	}

	return &entities.User{
		ID:        mongoModel.UserID,
		Email:     mongoModel.Email,
		Password:  mongoModel.Password,
		UserType:  valueobject.ParseUserType(mongoModel.UserType),
		CreatedAt: mongoModel.CreatedAt,
		UpdatedAt: mongoModel.UpdatedAt,
		DeletedAt: mongoModel.DeletedAt,
	}, nil
}

func (r *UserRepository) GetByEmail(email string) (*entities.User, error) {
	logger := logs.Logger()

	var mongoModel UserMongoDB
	filter := bson.M{"email": email, "deleted_at": bson.M{"$exists": false}}

	err := r.collection.FindOne(context.Background(), filter).Decode(&mongoModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logger.Error().Err(err).Str("email", email).Msg("Error finding user")
		return nil, err
	}

	return &entities.User{
		ID:        mongoModel.UserID,
		Email:     mongoModel.Email,
		Password:  mongoModel.Password,
		UserType:  valueobject.ParseUserType(mongoModel.UserType),
		CreatedAt: mongoModel.CreatedAt,
		UpdatedAt: mongoModel.UpdatedAt,
		DeletedAt: mongoModel.DeletedAt,
	}, nil
}

func (r *UserRepository) Create(user *entities.User) error {
	logger := logs.Logger()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(err).Msg("Error hashing password")
		return err
	}

	now := time.Now()
	mongoModel := &UserMongoDB{
		UserID:    user.ID,
		Email:     user.Email,
		Password:  string(hashedPassword),
		UserType:  user.UserType.String(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = r.collection.InsertOne(context.Background(), mongoModel)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating user in MongoDB")
		return err
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}
func (r *UserRepository) Update(user *entities.User) error {
	logger := logs.Logger()

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"email":      user.Email,
			"user_type":  user.UserType.String(),
			"updated_at": now,
		},
	}

	// Se a senha foi alterada, fazer hash
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error().Err(err).Msg("Error hashing password")
			return err
		}
		update["$set"].(bson.M)["password"] = string(hashedPassword)
	}

	filter := bson.M{"user_id": user.ID}
	result, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error updating user")
		return err
	}

	if result.MatchedCount == 0 {
		logger.Warn().Str("user_id", user.ID).Msg("User not found for update")
		return mongo.ErrNoDocuments
	}

	user.UpdatedAt = now
	return nil
}

func (r *UserRepository) Delete(id string) error {
	logger := logs.Logger()

	now := time.Now()
	filter := bson.M{"user_id": id, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"deleted_at": &now, "updated_at": now}}
	result, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		logger.Error().Err(err).Str("user_id", id).Msg("Error deleting user")
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *UserRepository) List() ([]entities.User, error) {
	logger := logs.Logger()

	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"deleted_at": bson.M{"$exists": false}},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		logger.Error().Err(err).Msg("Error listing users")
		return nil, err
	}
	defer cursor.Close(context.Background())

	result := make([]entities.User, 0)
	for cursor.Next(context.Background()) {
		var mongoModel UserMongoDB
		if err := cursor.Decode(&mongoModel); err != nil {
			logger.Error().Err(err).Msg("Error decoding user")
			continue
		}

		result = append(result, entities.User{
			ID:        mongoModel.UserID,
			Email:     mongoModel.Email,
			Password:  mongoModel.Password,
			UserType:  valueobject.ParseUserType(mongoModel.UserType),
			CreatedAt: mongoModel.CreatedAt,
			UpdatedAt: mongoModel.UpdatedAt,
			DeletedAt: mongoModel.DeletedAt,
		})
	}

	return result, nil
}
