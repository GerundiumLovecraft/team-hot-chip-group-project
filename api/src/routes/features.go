package routes

import (
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/controllers"
	"github.com/gin-gonic/gin"
)

func SetupFeaturesRouter(router *gin.RouterGroup) {
	features := router.Group("/features")

	features.GET("", controllers.GetAllFeatures)
}
