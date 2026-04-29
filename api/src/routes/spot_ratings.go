package routes

import (
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/controllers"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/middleware"
	"github.com/gin-gonic/gin"
)

func setupRatingRoutes(baseRouter *gin.RouterGroup) {
	spots := baseRouter.Group("/spots")
	spots.POST("/:id/rate", middleware.AuthenticationMiddleware, controllers.AddRating)
}