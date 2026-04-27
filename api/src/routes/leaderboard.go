package routes

import (
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/controllers"
	"github.com/gin-gonic/gin"
)

func setupLeaderboardRoutes(baseRouter *gin.RouterGroup) {
	leaderboard := baseRouter.Group("/leaderboard")
	leaderboard.GET("", controllers.GetLeaderboardData)
}