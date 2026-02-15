package routes

import (
	"strconv"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/gateway"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/billing_service"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/entity_api"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/execution_service"
	handlers "github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/middleware"
	repository "github.com/fiap-grupo95/os-service-api/internal/adapter/persistence"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/database"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/logs"
	"github.com/fiap-grupo95/os-service-api/internal/infrastructure/observability"
	"github.com/fiap-grupo95/os-service-api/internal/usecase"
	"github.com/fiap-grupo95/os-service-api/pkg/utils/auth"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var router = gin.Default()

const PORT = 8080

// Run will start the server
func Run() {
	newRelicApp, err := observability.NewRelicApp()
	if err != nil {
		logs.Logger().Error().Err(err).Msg("Failed to initialize New Relic; continuing without New Relic integration")
	}

	setMiddlewares()

	router.Use(nrgin.Middleware(newRelicApp))

	// Initialize logger with New Relic integration - Global Variable
	logs.Init(newRelicApp)

	// Initialize metrics collector with New Relic integration
	if newRelicApp != nil {
		observability.SetMetricsCollector(
			observability.NewNewRelicMetricsCollector(newRelicApp),
		)
	}

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	InitApp()

	logger := logs.Logger()
	err = router.Run(":" + strconv.Itoa(PORT))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}

func InitApp() {
	// Config JWT
	jwtCfg := auth.LoadJWTConfig()
	jwtService := auth.NewJWTService(jwtCfg)
	db := database.ConnectDatabase()

	// Handler de autenticação
	userRepository := repository.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(
		usecase.NewAuthUseCase(jwtService, userRepository),
	)

	// Auth routes
	addAuthRoutes(router, authHandler)

	partsSupplyRepository := entity_api.NewPartsSupplyRepository()
	serviceRepository := entity_api.NewServiceRepository()
	vehiclesRepository := entity_api.NewVehicleRepository()
	customerRepository := entity_api.NewCustomerRepository()
	serviceOrderRepository := repository.NewServiceOrderRepository(db)
	billingServiceRepository := billing_service.NewBillingServiceRepository()
	executionServiceRepository := execution_service.NewExecutionServiceRepository()
	// additionalRepairRepository := repository.NewAdditionalRepairRepository(db)

	vehiclesGateway := gateway.NewVehicleGateway(vehiclesRepository)
	customerGateway := gateway.NewCustomerGateway(customerRepository)
	partsSupplyGateway := gateway.NewPartsSupplyGateway(partsSupplyRepository)
	serviceGateway := gateway.NewServiceGateway(serviceRepository)
	billingServiceGateway := gateway.NewBillingServiceGateway(billingServiceRepository)
	executionGateway := gateway.NewExecutionServiceGateway(executionServiceRepository)
	serviceOrderGateway := gateway.NewServiceOrderGateway(
		serviceOrderRepository,
		vehiclesGateway,
		customerGateway,
		partsSupplyGateway,
		serviceGateway,
	)

	serviceOrderUsecase := usecase.NewServiceOrderUseCase(
		serviceOrderGateway,
		vehiclesGateway,
		customerGateway,
		serviceGateway,
		partsSupplyGateway,
		billingServiceGateway,
		executionGateway,
	)
	// additionalRepairUsecase := usecase.NewSOAdditionalRepairUseCase(
	// 	additionalRepairGateway,
	// 	serviceOrderGateway,
	// 	serviceGateway,
	// 	partsSupplyGateway)

	serviceOrderHandler := handlers.NewServiceOrderHandler(serviceOrderUsecase)
	// additionalRepairHandler := handlers.NewAdditionalRepairHandler(additionalRepairUsecase)

	// Protected routes
	router.Use(middleware.AuthMiddleware(jwtService))
	addPingRoutes(router)
	addServiceOrderRoutes(router, serviceOrderHandler)
	// addAdditionalRepairRoutes(router, additionalRepairHandler)
}

func setMiddlewares() {

	middleware.SetTrustedProxies(router)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logs.Logger().Error().Msgf("Recovered from panic: %v", recovered)
		c.AbortWithStatus(500)
	}))

}
