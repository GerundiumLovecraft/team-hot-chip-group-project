package routes

import (
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/controllers"
    "github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/middleware"
    "github.com/gin-gonic/gin"
)

func setupProfileRoutes(baseRouter *gin.RouterGroup) {
	profile := baseRouter.Group("/profile")

	profile.GET("", middleware.AuthenticationMiddleware, controllers.GetProfile)
	profile.GET("/spots", middleware.AuthenticationMiddleware, controllers.GetSpotsByUser)
}