package controllers

import (
	"net/http"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/auth"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/models"
	"github.com/gin-gonic/gin"
)

type JSONSpot struct {
	ID      uint   `json:"_id"`
	UserId uint `json:"user_id"`
	Name string `json:"name"`
	Address string `json:"address"`
	Description string `json:"description"`
	OpenFrom string `json:"open_from"`
	OpenTo string `json:"open_to"`
}

func GetAllSpots(ctx *gin.Context) {
	spots, errorMessage := models.FetchAllSpots()

	if errorMessage != nil {
		SendInternalError(ctx, errorMessage)
		return
	}

	value, _ := ctx.Get("userID")
	userID := value.(string)
	token, _ := auth.GenerateToken(userID)

	jsonSpots := make([]JSONSpot, 0)
	for _, spot := range *spots {
		jsonSpots = append(jsonSpots, JSONSpot{
			ID: spot.ID,
			UserId: spot.UserId,
			Name: spot.Name,
			Address: spot.Address,
			Description: spot.Description,
			OpenFrom: spot.OpenFrom,
			OpenTo: spot.OpenTo,
		})
	} 

	ctx.JSON(http.StatusOK, gin.H{"spots": jsonSpots, "token": token})


}

func GetSpotById(ctx *gin.Context) {
	spotId := ctx.Param("id")

	spotFound, errorMessage := models.FindSpot(spotId)

	if errorMessage != nil {
		SendInternalError(ctx, errorMessage)
		return
	}

	result := JSONSpot{
		ID: spotFound.ID,
		UserId: spotFound.UserId,
		Name: spotFound.Name,
		Address: spotFound.Address,
		Description: spotFound.Description,
		OpenFrom: spotFound.OpenFrom,
		OpenTo: spotFound.OpenTo,
	}

	value, _ := ctx.Get("userID")
	userID := value.(string)
	token, _ := auth.GenerateToken(userID)

	ctx.JSON(http.StatusOK, gin.H{"Spot": result, "token": token})

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

	if newSpot.UserId == 0 || newSpot.Name == "" || newSpot.Address == "" || newSpot.Description == "" || newSpot.OpenFrom == "" || newSpot.OpenTo == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The new spot's name, address, description, and opening hours must be supply"})
		return
	}
	
	_, errorMessage = newSpot.Save()
	if errorMessage != nil {
		SendInternalError(ctx, errorMessage)
		return
	}
	
	token, _ := auth.GenerateToken(userID)

	ctx.JSON(http.StatusCreated, gin.H{"message": "OK", "token": token})
}
