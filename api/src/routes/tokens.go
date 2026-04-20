package routes

import (
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/controllers"
	"github.com/gin-gonic/gin"
)

func setupAuthenticationRoutes(baseRouter *gin.RouterGroup) {
	tokensRouter := baseRouter.Group("/tokens")

	tokensRouter.POST("", controllers.CreateToken)
}
