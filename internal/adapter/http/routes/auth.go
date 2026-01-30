package routes

import (
	"mecanica_xpto/internal/adapter/http/handlers"
	"github.com/gin-gonic/gin"
)

func addAuthRoutes(rg *gin.Engine, authHandler *handlers.AuthHandler) {
	rg.POST(PostLogin, authHandler.Login)
}
