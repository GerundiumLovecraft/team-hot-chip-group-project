package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(engine *gin.Engine) {
	apiRouter := engine.Group("/")
	setupUserRoutes(apiRouter)
	setupAuthenticationRoutes(apiRouter)
	setupSpotsRoutes(apiRouter)
	SetupFeaturesRouter(apiRouter)
	setupProfileRoutes(apiRouter)
}
