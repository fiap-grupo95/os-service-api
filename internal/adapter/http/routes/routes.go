package routes

import (
	_ "mecanica_xpto/docs" // This will be auto-generated
	"mecanica_xpto/internal/adapter/http/handlers"
	"mecanica_xpto/internal/adapter/http/middleware"
	repository2 "mecanica_xpto/internal/adapter/persistence/repository"
	"mecanica_xpto/internal/infrastructure/database"
	"mecanica_xpto/internal/infrastructure/logs"
	"mecanica_xpto/internal/infrastructure/observability"
	"mecanica_xpto/internal/usecase"
	"mecanica_xpto/pkg/utils"
	"strconv"

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

	getRoutes()

	logger := logs.Logger()
	err = router.Run(":" + strconv.Itoa(PORT))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}

func getRoutes() {
	// Config JWT
	jwtCfg := utils.LoadJWTConfig()
	jwtService := utils.NewJWTService(jwtCfg)

	db := database.ConnectDatabase()
	userRepository := repository2.NewUserRepository(db)

	// Handler de autenticação
	authHandler := handlers.NewAuthHandler(
		usecase.NewAuthUseCase(jwtService, userRepository),
	)

	// Rotas PÚBLICAS
	v1 := router.Group("/v1")
	v1.POST("/login", authHandler.Login)

	partsSupplyRepository := repository2.NewPartsSupplyRepository(db)
	serviceRepository := repository2.NewServiceRepository(db)
	vehiclesRepository := repository2.NewVehicleRepository(db)
	customerRepository := repository2.NewCustomerRepository(db)
	serviceOrderRepository := repository2.NewServiceOrderRepository(db)
	additionalRepairRepository := repository2.NewAdditionalRepairRepository(db)

	serviceOrderUsecase := usecase.NewServiceOrderUseCase(
		serviceOrderRepository,
		vehiclesRepository,
		customerRepository,
		serviceRepository,
		partsSupplyRepository)
	additionalRepairUsecase := usecase.NewSOAdditionalRepairUseCase(
		additionalRepairRepository,
		serviceOrderRepository,
		serviceRepository,
		partsSupplyRepository)
		
	serviceOrderHandler := handlers.NewServiceOrderHandler(serviceOrderUsecase)
	additionalRepairHandler := handlers.NewAdditionalRepairHandler(additionalRepairUsecase)

	// Rotas PROTEGIDAS
	authGroup := v1.Group("/")
	authGroup.Use(middleware.AuthMiddleware(jwtService))
	addPingRoutes(authGroup)
	addServiceOrderRoutes(authGroup, serviceOrderHandler)
	addAdditionalRepairRoutes(authGroup, additionalRepairHandler)
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
