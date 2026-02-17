package routes

import (
	handlers "github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

func addAdditionalRepairRoutes(rg *gin.Engine, additionalRepair *handlers.AdditionalRepairHandler) {
		rg.POST(PostAdditionalRepair, additionalRepair.CreateAdditionalRepair)
		rg.GET(GetAdditionalRepair, additionalRepair.GetAdditionalRepair)
		rg.GET(GetAdditionalRepairBySO, additionalRepair.GetAdditionalRepairBySO)
		rg.POST(PostAdditionalRepairApprove, additionalRepair.ApproveAdditionalRepair)
		rg.POST(PostAdditionalRepairReject, additionalRepair.RejectAdditionalRepair)
		rg.POST(PostAdditionalRepairCancel, additionalRepair.CancelAdditionalRepair)
	}
