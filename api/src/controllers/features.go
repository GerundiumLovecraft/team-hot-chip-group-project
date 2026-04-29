package controllers

import (
	"net/http"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/models"
	"github.com/gin-gonic/gin"
)

type JSONListedFeature struct {
	ID   uint   `json:"feat_id"`
	Name string `json:"feat_name"`
}

func GetAllFeatures(ctx *gin.Context) {
	features, err := models.FetchAllFeatures()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	jsonFeatures := make([]JSONListedFeature, 0)

	for _, feat := range *features {
		jsonFeatures = append(jsonFeatures, JSONListedFeature{
			ID:   feat.ID,
			Name: feat.FeatName,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "List of features",
		"features": jsonFeatures,
	})
}