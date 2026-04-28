package controllers

import (
	"net/http"
	"strconv"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/auth"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/models"
	"github.com/gin-gonic/gin"
)

type JSONFeature struct {
	FeatName string `json:"feat_name"`
	Value    *int8  `json:"value"`
}

type JSONSpot struct {
	ID          uint          `json:"_id"`
	UserId      uint          `json:"user_id"`
	Name        string        `json:"name"`
	Address     string        `json:"address"`
	Description string        `json:"description"`
	Image       string        `json:"image"`
	OpenFrom    string        `json:"open_from"`
	OpenTo      string        `json:"open_to"`
	Features    []JSONFeature `json:"features"`
}

func GetAllSpots(ctx *gin.Context) {
	spots, errorMessage := models.FetchAllSpots()

	if errorMessage != nil {
		SendInternalError(ctx, errorMessage)
		return
	}

	jsonSpots := make([]JSONSpot, 0)
	for _, spot := range *spots {
		jsonFeature := make([]JSONFeature, 0)
		for _, feature := range spot.Features {
			jsonFeature = append(jsonFeature, JSONFeature{
				FeatName: feature.Feature.FeatName,
				Value:    feature.Value,
			})
		}
		jsonSpots = append(jsonSpots, JSONSpot{
			ID:          spot.ID,
			UserId:      spot.UserId,
			Name:        spot.Name,
			Address:     spot.Address,
			Description: spot.Description,
			Image:       spot.Image,
			OpenFrom:    spot.OpenFrom,
			OpenTo:      spot.OpenTo,
			Features:    jsonFeature,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"spots": jsonSpots})

}

func GetSpotById(ctx *gin.Context) {
	spotId := ctx.Param("id")

	spotFound, errorMessage := models.FindSpot(spotId)

	if errorMessage != nil {
		if errorMessage.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"spotId": JSONSpot{}})
			return
		}
		SendInternalError(ctx, errorMessage)
		return
	}

	jsonFeature := make([]JSONFeature, 0)
	for _, feature := range spotFound.Features {
		jsonFeature = append(jsonFeature, JSONFeature{
			FeatName: feature.Feature.FeatName,
			Value:    feature.Value,
		})
	}
	spotResult := JSONSpot{
		ID:          spotFound.ID,
		UserId:      spotFound.UserId,
		Name:        spotFound.Name,
		Address:     spotFound.Address,
		Description: spotFound.Description,
		Image:       spotFound.Image,
		OpenFrom:    spotFound.OpenFrom,
		OpenTo:      spotFound.OpenTo,
		Features:    jsonFeature,
	}

	ctx.JSON(http.StatusOK, gin.H{"Spot": spotResult})

}

func GetSpotsByFeature(ctx *gin.Context) {
	// get the slice of features from the request body
	features := make([]models.FeatureFilter, 0)
	if err := ctx.ShouldBindJSON(&features); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if features for filter were provided
	if len(features) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Filter parameters are empty"})
	}

	// filter spots by features (id and value)
	spotsFound, errorMessage := models.FilterSpotsByFeature(features)

	if errorMessage != nil {
		SendInternalError(ctx, errorMessage)
		return
	}

	// Return empty slice when no matches are found
	if len(*spotsFound) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Spots not found", "spots": []JSONSpot{}})
		return
	}

	jsonSpots := make([]JSONSpot, 0)
	for _, spot := range *spotsFound {
		jsonFeature := make([]JSONFeature, 0)
		for _, feature := range spot.Features {
			jsonFeature = append(jsonFeature, JSONFeature{
				FeatName: feature.Feature.FeatName,
				Value:    feature.Value,
			})
		}
		jsonSpots = append(jsonSpots, JSONSpot{
			ID:          spot.ID,
			UserId:      spot.UserId,
			Name:        spot.Name,
			Address:     spot.Address,
			Description: spot.Description,
			Image:       spot.Image,
			OpenFrom:    spot.OpenFrom,
			OpenTo:      spot.OpenTo,
			Features:    jsonFeature,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"spots": jsonSpots})

}

func GetSpotsByUser(ctx *gin.Context) {
    value, _ := ctx.Get("userID")
    userID := value.(string)

    spots, errorMessage := models.FetchSpotsByUserId(userID)
    if errorMessage != nil {
        SendInternalError(ctx, errorMessage)
        return
    }

    jsonSpots := make([]JSONSpot, 0)
    for _, spot := range *spots {
        jsonFeature := make([]JSONFeature, 0)
        for _, feature := range spot.Features {
            jsonFeature = append(jsonFeature, JSONFeature{
                FeatName: feature.Feature.FeatName,
                Value:    feature.Value,
            })
        }
        jsonSpots = append(jsonSpots, JSONSpot{
            ID:          spot.ID,
            UserId:      spot.UserId,
            Name:        spot.Name,
            Address:     spot.Address,
            Description: spot.Description,
            Image:       spot.Image,
            OpenFrom:    spot.OpenFrom,
            OpenTo:      spot.OpenTo,
            Features:    jsonFeature,
        })
    }

    ctx.JSON(http.StatusOK, gin.H{"spots": jsonSpots})
}

func CreateSpot(ctx *gin.Context) {
	var newSpot models.Spot
	errorMessage := ctx.BindJSON(&newSpot)

	if errorMessage != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": errorMessage.Error()})
		return
	}
	value, _ := ctx.Get("userID")
	userID := value.(string)

	userIdUint, _ := strconv.ParseUint(userID, 10, 64)

	newSpot.UserId = uint(userIdUint)

	if newSpot.UserId == 0 || newSpot.Name == "" || newSpot.Address == "" || newSpot.Description == "" || newSpot.OpenFrom == "" || newSpot.OpenTo == "" || len(newSpot.Features) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The new spot's name, address, description, opening hours, and features must be supply"})
		return
	}

	_, errorMessage = newSpot.Save()
	if errorMessage != nil {
		SendInternalError(ctx, errorMessage)
		return
	}

	token, _ := auth.GenerateToken(userID)

	ctx.JSON(http.StatusCreated, gin.H{"spotID": newSpot.ID, "token": token})
}
