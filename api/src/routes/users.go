package routes

import (
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/controllers"
	"github.com/gin-gonic/gin"
)

func setupUserRoutes(baseRouter *gin.RouterGroup) {
	users := baseRouter.Group("/users")

	users.POST("", controllers.CreateUser)
}
