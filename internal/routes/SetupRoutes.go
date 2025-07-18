package routes

import (
	"github.com/gin-gonic/gin"
	"note1/internal/services"
)

func SetupRoutes(r *gin.Engine, serviceContainer *services.ServicesContainer) {
	noteRoutesSetup(r, serviceContainer)
	userRoutesSetup(r, serviceContainer)
}
