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
	ID                 string  `bson:"id" json:"id"`
	Value              float64 `bson:"value" json:"value"`
	ServiceOrderID     string  `bson:"service_order_id" json:"service_order_id"`
	AdditionalRepairID string  `bson:"additional_repair_id" json:"additional_repair_id"`
	Status             string  `bson:"status" json:"status"`
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
		ServiceOrderID: additionalRepair.ServiceOrderID,
		Description:    additionalRepair.Description,
		Status:         additionalRepair.Status.String(),
		Estimate:       estimate,
		PartsSupplies:  partsSupplies,
		Services:       services,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	result, err := r.collection.InsertOne(ctx, mongoModel)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating additional repair in MongoDB")
		return nil, err
	}

	mongoModel.ID = result.InsertedID.(primitive.ObjectID)
	additionalRepair.ID = mongoModel.ID.Hex()
	additionalRepair.CreatedAt = now
	additionalRepair.UpdatedAt = now
	return additionalRepair, nil
}

func (r *AdditionalRepairRepository) GetByID(ctx context.Context, id string) (*entities.AdditionalRepair, error) {
	logger := logs.LoggerWithContext(ctx)

	var mongoModel AdditionalRepairMongoDB
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}
	filter := bson.M{"_id": objID}

	err = r.collection.FindOne(ctx, filter).Decode(&mongoModel)
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
		ID:             mongoModel.ID.Hex(),
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

func (r *AdditionalRepairRepository) GetByServiceOrderID(ctx context.Context, serviceOrderID string) ([]entities.AdditionalRepair, error) {
	logger := logs.LoggerWithContext(ctx)

	filter := bson.M{"service_order_id": serviceOrderID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		logger.Error().Err(err).Str("service_order_id", serviceOrderID).Msg("Error listing additional repairs by service order id")
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []entities.AdditionalRepair
	for cursor.Next(ctx) {
		var mongoModel AdditionalRepairMongoDB
		if err := cursor.Decode(&mongoModel); err != nil {
			logger.Error().Err(err).Str("service_order_id", serviceOrderID).Msg("Error decoding additional repair")
			return nil, err
		}

		var estimate *entities.Estimate
		if mongoModel.Estimate != nil {
			estimate = &entities.Estimate{
				ID:                 mongoModel.Estimate.ID,
				ServiceOrderID:     mongoModel.Estimate.ServiceOrderID,
				AdditionalRepairID: mongoModel.Estimate.AdditionalRepairID,
				Value:              mongoModel.Estimate.Value,
				Status:             mongoModel.Estimate.Status,
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

		results = append(results, entities.AdditionalRepair{
			ID:             mongoModel.ID.Hex(),
			Description:    mongoModel.Description,
			ServiceOrderID: mongoModel.ServiceOrderID,
			Status:         valueobject.ParseAdditionalRepairStatus(mongoModel.Status),
			Estimate:       estimate,
			CreatedAt:      mongoModel.CreatedAt,
			UpdatedAt:      mongoModel.UpdatedAt,
			PartsSupplies:  partsSupplies,
			Services:       services,
		})
	}
	if err := cursor.Err(); err != nil {
		logger.Error().Err(err).Str("service_order_id", serviceOrderID).Msg("Error iterating additional repairs")
		return nil, err
	}

	if results == nil {
		results = []entities.AdditionalRepair{}
	}
	return results, nil
}

func (r *AdditionalRepairRepository) UpdateAdditionalRepair(ctx context.Context, additionalRepair *entities.AdditionalRepair) (*entities.AdditionalRepair, error) {
	logger := logs.LoggerWithContext(ctx)

	now := time.Now()
	if additionalRepair == nil || additionalRepair.ID == "" {
		return nil, mongo.ErrNoDocuments
	}

	objID, err := primitive.ObjectIDFromHex(additionalRepair.ID)
	if err != nil {
		return nil, mongo.ErrNoDocuments
	}

	set := bson.M{
		"updated_at": now,
	}
	if additionalRepair.Description != "" {
		set["description"] = additionalRepair.Description
	}
	if additionalRepair.ServiceOrderID != "" {
		set["service_order_id"] = additionalRepair.ServiceOrderID
	}
	if additionalRepair.Status.String() != "" {
		set["status"] = additionalRepair.Status.String()
	}
	if additionalRepair.Estimate != nil {
		if additionalRepair.Estimate.ID != "" {
			set["estimate.id"] = additionalRepair.Estimate.ID
		}
		if additionalRepair.Estimate.Status != "" {
			set["estimate.status"] = additionalRepair.Estimate.Status
		}
		if additionalRepair.Estimate.Value != 0 {
			set["estimate.value"] = additionalRepair.Estimate.Value
		}
		if additionalRepair.Estimate.ServiceOrderID != "" {
			set["estimate.service_order_id"] = additionalRepair.Estimate.ServiceOrderID
		}
		if additionalRepair.Estimate.AdditionalRepairID != "" {
			set["estimate.additional_repair_id"] = additionalRepair.Estimate.AdditionalRepairID
		}
	}

	update := bson.M{"$set": set}
	filter := bson.M{"_id": objID}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated AdditionalRepairMongoDB
	err = r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		logger.Error().Err(err).Str("ar_id", additionalRepair.ID).Msg("Error updating additional repair")
		return nil, err
	}

	var estimate *entities.Estimate
	if updated.Estimate != nil {
		estimate = &entities.Estimate{
			ID:                 updated.Estimate.ID,
			ServiceOrderID:     updated.Estimate.ServiceOrderID,
			AdditionalRepairID: updated.Estimate.AdditionalRepairID,
			Value:              updated.Estimate.Value,
			Status:             updated.Estimate.Status,
		}
	}

	partsSupplies := make([]entities.PartsSupply, 0, len(updated.PartsSupplies))
	for _, ps := range updated.PartsSupplies {
		partsSupplies = append(partsSupplies, entities.PartsSupply{
			ID:       ps.ID,
			Quantity: ps.Quantity,
			Price:    ps.Price,
		})
	}

	services := make([]entities.Service, 0, len(updated.Services))
	for _, s := range updated.Services {
		services = append(services, entities.Service{
			ID:    s.ID,
			Price: s.Price,
		})
	}

	return &entities.AdditionalRepair{
		ID:             updated.ID.Hex(),
		Description:    updated.Description,
		ServiceOrderID: updated.ServiceOrderID,
		Status:         valueobject.ParseAdditionalRepairStatus(updated.Status),
		Estimate:       estimate,
		CreatedAt:      updated.CreatedAt,
		UpdatedAt:      updated.UpdatedAt,
		PartsSupplies:  partsSupplies,
		Services:       services,
	}, nil
}

func (r *AdditionalRepairRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
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
