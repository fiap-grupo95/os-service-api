package routes

import (
	handlers "github.com/fiap-grupo95/os-service-api/internal/adapter/http/handlers"
	"github.com/gin-gonic/gin"
)

func addAuthRoutes(rg *gin.Engine, authHandler *handlers.AuthHandler) {
	rg.POST(PostLogin, authHandler.Login)
}
