package routes

import (
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/controllers"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/middleware"
	"github.com/gin-gonic/gin"
)

func setupUserRoutes(baseRouter *gin.RouterGroup) {
	users := baseRouter.Group("/users")

	users.POST("", controllers.CreateUser)
	users.GET("/search-by-username/:username", middleware.AuthenticationMiddleware, controllers.GetUserByUsername)
	users.GET("/:id", middleware.AuthenticationMiddleware, controllers.GetUserById)

	//Optional routes for getting a profile, editing profile and deleting
	//If we get to that point, we'll need to implement functions in respective files

	//Optional edit user
	//users.PUT("/:id", controllers.UpdateUser)

	//Option delete user
	//users.DELETE("/:id", controllers.DeleteUser)

}
