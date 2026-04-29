package controllers

import (
	"net/http"
	"strconv"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/models"
	"github.com/gin-gonic/gin"
)

// get the user id from the token
func AddRating(ctx *gin.Context) {
	userId, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorised"})
		return
	}

	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID"})
		return
	}

	// get the spot id from the url
	spotIdStr := ctx.Param("id")
	spotIdUint, err := strconv.ParseUint(spotIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid spot ID"})
		return
	}

	// get the rating from the request body
	var requestBody struct {
		Rating int8 `json:"rating"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// check rating is a valid number between 1 and 5
	if requestBody.Rating < 1 || requestBody.Rating > 5 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Rating must be between 1 and 5"})
		return
	}

	userIdUint, _ := strconv.ParseUint(userIdStr, 10, 64)

	// save the rating
	err = models.AddRating(uint(userIdUint), uint(spotIdUint), requestBody.Rating)
	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Rating added successfully"})
}
