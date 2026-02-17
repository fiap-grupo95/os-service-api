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

type ServiceOrderMongoDB struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CustomerID    string             `bson:"customer_id" json:"customer_id"`
	VehicleID     string             `bson:"vehicle_id" json:"vehicle_id"`
	Status        string             `bson:"status" json:"status"`
	Estimate      *EstimateMongoDB   `bson:"estimate,omitempty" json:"estimate,omitempty"`
	Execution     *ExecutionMongoDB  `bson:"execution,omitempty" json:"execution,omitempty"`
	PartsSupplies []PartsSupplyItem  `bson:"parts_supplies,omitempty" json:"parts_supplies,omitempty"`
	Services      []ServiceItem      `bson:"services,omitempty" json:"services,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type ExecutionMongoDB struct {
	ID             string    `bson:"id" json:"id"`
	ServiceOrderID string    `bson:"service_order_id" json:"service_order_id"`
	Status         string    `bson:"status" json:"status"`
	StartedAt      *time.Time `bson:"started_at" json:"started_at"`
	FinishedAt     *time.Time `bson:"finished_at" json:"finished_at"`
}

type ServiceOrderRepository struct {
	collection *mongo.Collection
}

var _ interfaces.IServiceOrderRepository = (*ServiceOrderRepository)(nil)

func NewServiceOrderRepository(db *database.MongoDBConfig) *ServiceOrderRepository {
	repo := &ServiceOrderRepository{
		collection: db.GetCollection("service_orders"),
	}
	_ = repo.CreateIndexes(context.Background())
	return repo
}

func (r *ServiceOrderRepository) Create(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	logger := logs.Logger()

	now := time.Now()

	mongoModel := &ServiceOrderMongoDB{
		CustomerID: serviceOrder.CustomerID,
		VehicleID:  serviceOrder.VehicleID,
		Status:     serviceOrder.Status.String(),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	result, err := r.collection.InsertOne(ctx, mongoModel)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating service order in MongoDB")
		return nil, err
	}

	mongoModel.ID = result.InsertedID.(primitive.ObjectID)

	serviceOrder.ID = mongoModel.ID.Hex()
	serviceOrder.CreatedAt = &now
	serviceOrder.UpdatedAt = &now
	return serviceOrder, nil
}

func (r *ServiceOrderRepository) GetByID(ctx context.Context, osId string) (*entities.ServiceOrder, error) {
	logger := logs.Logger()
	objectID, err := primitive.ObjectIDFromHex(osId)
	if err != nil {
		logger.Warn().Str("_id", osId).Msg("Invalid MongoDB ObjectID")
		return nil, nil
	}

	var mongoModel ServiceOrderMongoDB
	filter := bson.M{"_id": objectID}

	err = r.collection.FindOne(ctx, filter).Decode(&mongoModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Warn().Str("_id", osId).Msg("Service order not found")
			return nil, nil
		}
		logger.Error().Err(err).Str("_id", osId).Msg("Error finding service order")
		return nil, err
	}

	serviceOrder := &entities.ServiceOrder{
		ID:         mongoModel.ID.Hex(),
		CustomerID: mongoModel.CustomerID,
		VehicleID:  mongoModel.VehicleID,
		Status:     valueobject.ParseServiceOrderStatus(mongoModel.Status),
	}

	if mongoModel.Estimate != nil {
		estimate := entities.Estimate{
			ID:     mongoModel.Estimate.ID,
			Value:  mongoModel.Estimate.Value,
			Status: mongoModel.Estimate.Status,
		}
		serviceOrder.Estimate = &estimate
	}

	if mongoModel.Execution != nil {
		serviceOrder.Execution = &entities.Execution{
			ID:             mongoModel.Execution.ID,
			ServiceOrderID: mongoModel.Execution.ServiceOrderID,
			Status:         mongoModel.Execution.Status,
			StartedAt:      mongoModel.Execution.StartedAt,
			FinishedAt:     mongoModel.Execution.FinishedAt,
		}
	}

	serviceOrder.CreatedAt = &mongoModel.CreatedAt
	serviceOrder.UpdatedAt = &mongoModel.UpdatedAt

	return serviceOrder, nil
}

func (r *ServiceOrderRepository) Update(ctx context.Context, serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	logger := logs.Logger()
	objectID, err := primitive.ObjectIDFromHex(serviceOrder.ID)
	if err != nil {
		logger.Warn().Str("_id", serviceOrder.ID).Msg("Invalid MongoDB ObjectID")
		return nil, mongo.ErrNoDocuments
	}

	now := time.Now()

	setFields := bson.M{
		"updated_at": now,
	}

	if serviceOrder.CustomerID != "" {
		setFields["customer_id"] = serviceOrder.CustomerID
	}
	if serviceOrder.VehicleID != "" {
		setFields["vehicle_id"] = serviceOrder.VehicleID
	}
	if serviceOrder.Status.IsValid() {
		setFields["status"] = serviceOrder.Status.String()
	}
	if serviceOrder.Estimate != nil {
		v := serviceOrder.Estimate.Value
		setFields["estimate"] = &EstimateMongoDB{
			ID:             serviceOrder.Estimate.ID,
			Value:          v,
			ServiceOrderID: serviceOrder.ID,
			Status:         serviceOrder.Estimate.Status,
		}
	}

	if serviceOrder.Execution != nil{
		setFields["execution"] = &ExecutionMongoDB{
			ID:             serviceOrder.Execution.ID,
			ServiceOrderID: serviceOrder.ID,
			Status:         serviceOrder.Execution.Status,
			StartedAt:      serviceOrder.Execution.StartedAt,
			FinishedAt:     serviceOrder.Execution.FinishedAt,
		}
	}
	if serviceOrder.PartsSupplies != nil {
		partsSupplies := make([]bson.M, 0, len(serviceOrder.PartsSupplies))
		for _, ps := range serviceOrder.PartsSupplies {
			item := bson.M{
				"id":       ps.ID,
				"quantity": ps.Quantity,
			}
			if ps.Price != 0 {
				item["price"] = ps.Price
			}
			partsSupplies = append(partsSupplies, item)
		}
		setFields["parts_supplies"] = partsSupplies
	}
	if serviceOrder.Services != nil {
		services := make([]bson.M, 0, len(serviceOrder.Services))
		for _, s := range serviceOrder.Services {
			item := bson.M{
				"id": s.ID,
			}
			if s.Price != 0 {
				item["price"] = s.Price
			}
			services = append(services, item)
		}
		setFields["services"] = services
	}

	update := bson.M{"$set": setFields}

	filter := bson.M{"_id": objectID}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error().Err(err).Str("_id", serviceOrder.ID).Msg("Error updating service order")
		return nil, err
	}

	if result.MatchedCount == 0 {
		logger.Warn().Str("_id", serviceOrder.ID).Msg("Service order not found for update")
		return nil, mongo.ErrNoDocuments
	}

	serviceOrder.UpdatedAt = &now
	return serviceOrder, nil
}

func (r *ServiceOrderRepository) List(ctx context.Context) ([]*entities.ServiceOrder, error) {
	logger := logs.Logger()

	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		logger.Error().Err(err).Msg("Error listing service orders")
		return nil, err
	}
	defer cursor.Close(ctx)

	var serviceOrders []*entities.ServiceOrder
	for cursor.Next(ctx) {
		var mongoModel ServiceOrderMongoDB
		if err := cursor.Decode(&mongoModel); err != nil {
			logger.Error().Err(err).Msg("Error decoding service order")
			continue
		}

		createdAt := mongoModel.CreatedAt
		updatedAt := mongoModel.UpdatedAt

		serviceOrder := &entities.ServiceOrder{
			ID:         mongoModel.ID.Hex(),
			CustomerID: mongoModel.CustomerID,
			VehicleID:  mongoModel.VehicleID,
			Status:     valueobject.ParseServiceOrderStatus(mongoModel.Status),
			CreatedAt:  &createdAt,
			UpdatedAt:  &updatedAt,
		}
		if mongoModel.Estimate != nil {
			estimate := &entities.Estimate{
				ID:     mongoModel.Estimate.ID,
				Value:  mongoModel.Estimate.Value,
				Status: mongoModel.Estimate.Status,
			}
			serviceOrder.Estimate = estimate
		}

		if mongoModel.Execution != nil{
			execution := &entities.Execution{
				ID: mongoModel.Execution.ID,
				ServiceOrderID: mongoModel.Execution.ServiceOrderID,
				Status: mongoModel.Execution.Status,
				StartedAt: mongoModel.Execution.StartedAt,
				FinishedAt: mongoModel.Execution.FinishedAt,
			}
			serviceOrder.Execution = execution
		}

		serviceOrders = append(serviceOrders, serviceOrder)
	}

	return serviceOrders, nil
}

func (r *ServiceOrderRepository) CreateIndexes(ctx context.Context) error {
	logger := logs.Logger()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "customer_id", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "vehicle_id", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating indexes")
		return err
	}

	logger.Info().Msg("MongoDB indexes created successfully")
	return nil
}
