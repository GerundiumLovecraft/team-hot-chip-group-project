package main

import (

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/env"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/models"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	env.LoadEnv()

	app := setupApp()

	models.OpenDatabaseConnection()
	models.AutoMigrateModels()
	models.SeedFeatures()
	models.SeedDemoData()

	app.Run(":8082")
}

func setupApp() *gin.Engine {
	app := gin.Default()
	setupCORS(app)
	routes.SetupRoutes(app)
	return app
}

func setupCORS(app *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"}

	app.Use(cors.New(config))
}
