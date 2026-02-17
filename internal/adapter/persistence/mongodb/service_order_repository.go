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
	PartsSupplies []PartsSupplyItem  `bson:"parts_supplies,omitempty" json:"parts_supplies,omitempty"`
	Services      []ServiceItem      `bson:"services,omitempty" json:"services,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
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

	createdAt := mongoModel.CreatedAt
	updatedAt := mongoModel.UpdatedAt
	return &entities.ServiceOrder{
		ID:         mongoModel.ID.Hex(),
		CustomerID: mongoModel.CustomerID,
		VehicleID:  mongoModel.VehicleID,
		Status:     valueobject.ParseServiceOrderStatus(mongoModel.Status),
		CreatedAt:  &createdAt,
		UpdatedAt:  &updatedAt,
	}, nil
}

func (r *ServiceOrderRepository) Update(ctx context.Context, serviceOrder *entities.ServiceOrder) error {
	logger := logs.Logger()
	objectID, err := primitive.ObjectIDFromHex(serviceOrder.ID)
	if err != nil {
		logger.Warn().Str("_id", serviceOrder.ID).Msg("Invalid MongoDB ObjectID")
		return mongo.ErrNoDocuments
	}

	now := time.Now()

	// Atualizar campos no MongoDB
	update := bson.M{
		"$set": bson.M{
			"customer_id": serviceOrder.CustomerID,
			"vehicle_id":  serviceOrder.VehicleID,
			"status":      serviceOrder.Status.String(),
			"updated_at":  now,
		},
	}

	filter := bson.M{"_id": objectID}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error().Err(err).Str("_id", serviceOrder.ID).Msg("Error updating service order")
		return err
	}

	if result.MatchedCount == 0 {
		logger.Warn().Str("_id", serviceOrder.ID).Msg("Service order not found for update")
		return mongo.ErrNoDocuments
	}

	serviceOrder.UpdatedAt = &now
	return nil
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
		serviceOrders = append(serviceOrders, &entities.ServiceOrder{
			ID:         mongoModel.ID.Hex(),
			CustomerID: mongoModel.CustomerID,
			VehicleID:  mongoModel.VehicleID,
			Status:     valueobject.ParseServiceOrderStatus(mongoModel.Status),
			CreatedAt:  &createdAt,
			UpdatedAt:  &updatedAt,
		})
	}

	return serviceOrders, nil
}

func (r *ServiceOrderRepository) Delete(ctx context.Context, osId string) error {
	logger := logs.Logger()
	objectID, err := primitive.ObjectIDFromHex(osId)
	if err != nil {
		logger.Warn().Str("_id", osId).Msg("Invalid MongoDB ObjectID")
		return mongo.ErrNoDocuments
	}

	filter := bson.M{"_id": objectID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		logger.Error().Err(err).Str("_id", osId).Msg("Error deleting service order")
		return err
	}

	if result.DeletedCount == 0 {
		logger.Warn().Str("_id", osId).Msg("Service order not found for deletion")
		return mongo.ErrNoDocuments
	}

	return nil
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
