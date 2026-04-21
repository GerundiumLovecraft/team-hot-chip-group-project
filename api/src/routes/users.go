package routes

import (
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/controllers"
	"github.com/gin-gonic/gin"
)

func setupUserRoutes(baseRouter *gin.RouterGroup) {
	users := baseRouter.Group("/users")

	//Empty string path "" is treated as users in Gin due to it being users.POST
	users.POST("", controllers.CreateUser)

	//Optional routes for getting a profile, editing profile and deleting
	//If we get to that point, we'll need to implement functions in respective files

	//Optional get user profile
	// users.GET("/:id", controllers.GetUser)

	//Optional edit user
	//users.PUT("/:id", controllers.UpdateUser)

	//Option delete user
	//users.DELETE("/:id", controllers.DeleteUser)
}
