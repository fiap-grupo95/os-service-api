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
)

type AdditionalRepairMongoDB struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ARId           string             `bson:"ar_id" json:"ar_id"`
	ServiceOrderID string             `bson:"service_order_id" json:"service_order_id"`
	Description    string             `bson:"description" json:"description"`
	Status         string             `bson:"status" json:"status"`
	Estimate       *EstimateMongoDB   `bson:"estimate,omitempty" json:"estimate,omitempty"`
	PartsSupplies  []PartsSupplyItem  `bson:"parts_supplies,omitempty" json:"parts_supplies,omitempty"`
	Services       []ServiceItem      `bson:"services,omitempty" json:"services,omitempty"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type EstimateMongoDB struct {
	ID                 string   `bson:"id" json:"id"`
	Value              float64 `bson:"value" json:"value"`
	ServiceOrderID     string   `bson:"service_order_id" json:"service_order_id"`
	AdditionalRepairID string   `bson:"additional_repair_id" json:"additional_repair_id"`
	Status             string   `bson:"status" json:"status"`
}

type PartsSupplyItem struct {
	ID       string  `bson:"id" json:"id"`
	Quantity int     `bson:"quantity" json:"quantity"`
	Price    float64 `bson:"price" json:"price"`
}

type ServiceItem struct {
	ID    string  `bson:"id" json:"id"`
	Price float64 `bson:"price" json:"price"`
}

type AdditionalRepairRepository struct {
	collection *mongo.Collection
}

var _ interfaces.IAdditionalRepairRepository = (*AdditionalRepairRepository)(nil)

func NewAdditionalRepairRepository(db *database.MongoDBConfig) *AdditionalRepairRepository {
	repo := &AdditionalRepairRepository{
		collection: db.GetCollection("additional_repairs"),
	}
	_ = repo.CreateIndexes(context.Background())
	return repo
}

func (r *AdditionalRepairRepository) CreateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error) {
	logger := logs.LoggerWithContext(ctx)

	now := time.Now()

	var estimate *EstimateMongoDB
	if additionalRepair.Estimate != nil {
		v := additionalRepair.Estimate.Value
		estimate = &EstimateMongoDB{
			ID:                 additionalRepair.Estimate.ID,
			Value:              v,
			AdditionalRepairID: additionalRepair.ID,
			Status:             additionalRepair.Estimate.Status,
		}
	}

	partsSupplies := make([]PartsSupplyItem, 0, len(additionalRepair.PartsSupplies))
	for _, ps := range additionalRepair.PartsSupplies {
		partsSupplies = append(partsSupplies, PartsSupplyItem{
			ID:       ps.ID,
			Quantity: ps.Quantity,
			Price:    ps.Price,
		})
	}

	services := make([]ServiceItem, 0, len(additionalRepair.Services))
	for _, s := range additionalRepair.Services {
		services = append(services, ServiceItem{
			ID:    s.ID,
			Price: s.Price,
		})
	}

	mongoModel := &AdditionalRepairMongoDB{
		ARId:           additionalRepair.ID,
		ServiceOrderID: additionalRepair.ServiceOrderID,
		Description:    additionalRepair.Description,
		Status:         additionalRepair.Status.String(),
		Estimate:       estimate,
		PartsSupplies:  partsSupplies,
		Services:       services,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	_, err := r.collection.InsertOne(ctx, mongoModel)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating additional repair in MongoDB")
		return nil, err
	}

	additionalRepair.CreatedAt = now
	additionalRepair.UpdatedAt = now
	return additionalRepair, nil
}

func (r *AdditionalRepairRepository) GetByID(ctx context.Context, id string) (*entities.AdditionalRepair, error) {
	logger := logs.LoggerWithContext(ctx)

	var mongoModel AdditionalRepairMongoDB
	filter := bson.M{"ar_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&mongoModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logger.Error().Err(err).Str("ar_id", id).Msg("Error finding additional repair")
		return nil, err
	}

	var estimate *entities.Estimate
	if mongoModel.Estimate != nil {
		estimate = &entities.Estimate{
			ID:     mongoModel.Estimate.ID,
			Value:  mongoModel.Estimate.Value,
			Status: mongoModel.Estimate.Status,
		}
	}

	partsSupplies := make([]entities.PartsSupply, 0, len(mongoModel.PartsSupplies))
	for _, ps := range mongoModel.PartsSupplies {
		partsSupplies = append(partsSupplies, entities.PartsSupply{
			ID:       ps.ID,
			Quantity: ps.Quantity,
			Price:    ps.Price,
		})
	}

	services := make([]entities.Service, 0, len(mongoModel.Services))
	for _, s := range mongoModel.Services {
		services = append(services, entities.Service{
			ID:    s.ID,
			Price: s.Price,
		})
	}

	return &entities.AdditionalRepair{
		ID:             mongoModel.ARId,
		Description:    mongoModel.Description,
		ServiceOrderID: mongoModel.ServiceOrderID,
		Status:         valueobject.ParseAdditionalRepairStatus(mongoModel.Status),
		Estimate:       estimate,
		CreatedAt:      mongoModel.CreatedAt,
		UpdatedAt:      mongoModel.UpdatedAt,
		PartsSupplies:  partsSupplies,
		Services:       services,
	}, nil
}

func (r *AdditionalRepairRepository) UpdateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error) {
	logger := logs.LoggerWithContext(ctx)

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"description":      additionalRepair.Description,
			"status":           additionalRepair.Status.String(),
			"service_order_id": additionalRepair.ServiceOrderID,
			"updated_at":       now,
		},
	}

	filter := bson.M{"ar_id": additionalRepair.ID}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error().Err(err).Str("ar_id", additionalRepair.ID).Msg("Error updating additional repair")
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	additionalRepair.UpdatedAt = now
	return additionalRepair, nil
}

func (r *AdditionalRepairRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "ar_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "service_order_id", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: -1}},
		},
	}
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
