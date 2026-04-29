package controllers

import (
	"net/http"
	"strconv"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/auth"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/models"
	"github.com/gin-gonic/gin"
)

// GetProfile returns the profile of the currently authenticated user.
// It relies on the AuthenticationMiddleware having already set "userID" in the context.
// Route: GET /profile  (protected by AuthenticationMiddleware)

func GetProfile(ctx *gin.Context) {

	userId, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorised"})
		return
	}

	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID in token"})
		return
	}

	user, err := models.FindUser(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	token, err := auth.GenerateToken(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	userIdAsStr := strconv.FormatUint(uint64(user.ID), 10)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"token":   token,
		"user": gin.H{
			"id":        userIdAsStr,
			"username":  user.Username,
			"email":     user.Email,
			"createdAt": user.CreatedAt,
			"avatar": user.Avatar,
		},
	})
}