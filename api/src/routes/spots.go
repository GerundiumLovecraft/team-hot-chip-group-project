package routes

import (
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/controllers"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/middleware"
	"github.com/gin-gonic/gin"
)

func setupSpotsRoutes(baseRouter *gin.RouterGroup) {
	spots := baseRouter.Group("/spots")

	spots.GET("", controllers.GetAllSpots)
	spots.GET("/:id", controllers.GetSpotById)
	spots.POST("", middleware.AuthenticationMiddleware, controllers.CreateSpot)

}
