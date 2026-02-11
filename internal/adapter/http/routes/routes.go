package routes

import (
	"net/http"
	"strconv"

	handlers "github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers"
	middleware "github.com/fiap-grupo95/os-service-api/internal/adapter/http/middleware"
	entityapi "github.com/fiap-grupo95/os-service-api/internal/adapter/persistence/entity_api"
	"github.com/fiap-grupo95/os-service-api/internal/adapter/persistence/gateway"
	repository "github.com/fiap-grupo95/os-service-api/internal/adapter/persistence/repository"
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

	partsSupplyRepository := repository.NewPartsSupplyRepository(db)
	serviceRepository := repository.NewServiceRepository(db)
	vehiclesRepository := entityapi.NewVehicleRepository()
	customerRepository := entityapi.NewCustomerRepository(&http.Client{})
	serviceOrderRepository := repository.NewServiceOrderRepository(db)
	// additionalRepairRepository := repository.NewAdditionalRepairRepository(db)

	serviceOrderGateway := gateway.NewServiceOrderGateway(serviceOrderRepository)
	vehiclesGateway := gateway.NewVehicleGateway(vehiclesRepository)
	customerGateway := gateway.NewCustomerGateway(customerRepository)

	serviceOrderUsecase := usecase.NewServiceOrderUseCase(
		serviceOrderGateway,
		vehiclesGateway,
		customerGateway,
		serviceRepository,
		partsSupplyRepository)
	// additionalRepairUsecase := usecase.NewSOAdditionalRepairUseCase(
	// 	additionalRepairGateway,
	// 	serviceOrderGateway,
	// 	serviceRepository,
	// 	partsSupplyRepository)

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
