package routes

import (
	handlers "github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

func addAdditionalRepairRoutes(rg *gin.Engine, additionalRepair *handlers.AdditionalRepairHandler) {
	serviceOrdersRoutes := rg.Group(PathAdditionalRepair)
	{
		serviceOrdersRoutes.POST("", additionalRepair.CreateAdditionalRepair)
		serviceOrdersRoutes.GET("/:id", additionalRepair.GetAdditionalRepair)
		serviceOrdersRoutes.PATCH("/:id/add", additionalRepair.AddPartSupplyAndService)
		serviceOrdersRoutes.PATCH("/:id/remove", additionalRepair.RemovePartSupplyAndService)
		serviceOrdersRoutes.PATCH("/:id/customer_approval", additionalRepair.CustomerApproval)
	}
}
