package controllers

import (
	"net/http"
	"strconv"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/auth"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(ctx *gin.Context) {
	var newUser models.User
	err := ctx.BindJSON(&newUser)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if newUser.Email == "" || newUser.HashedPassword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Must supply username and password"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		SendInternalError(ctx, err)
		return
	}
	newUser.HashedPassword = string(hashedPassword)

	_, err = newUser.Save()
	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "OK"})
}

func GetUserById(ctx *gin.Context) {
	// Get the user ID from the token
	userId, exists := ctx.Get("userID")
	userIdStr, ok := userId.(string)

	// Get the user ID from the parameters
	searchedId := ctx.Param("id")

	// Check if the user's ID matches to the searched ID
	if !exists || !ok || userIdStr != searchedId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorised access to user's profile"})
		return
	}

	// Find the use by ID
	user, err := models.FindUser(searchedId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User with this ID not found"})
		return
	}

	// Generate new token
	token, _ := auth.GenerateToken(userIdStr)

	// Send the response
	ctx.JSON(http.StatusOK, gin.H{"message": "OK", "user": user, "token": token})
}

func GetUserByUsername(ctx *gin.Context) {
	// Search for the user by the username
	username := ctx.Param("username")
	user, err := models.FindUserByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Username not found"})
		return
	}

	// To keep the page private for now, check if the user's ID from the token matches with the ID in user object
	userId, exists := ctx.Get("userID")
	userIdStr, ok := userId.(string)
	uIntToStr := strconv.FormatUint(uint64(user.ID), 10)
	if !exists || !ok || userIdStr != uIntToStr {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorised access to user's profile"})
		return
	}
	// Generate new token
	token, _ := auth.GenerateToken(userIdStr)
	// Send the response
	ctx.JSON(http.StatusOK, gin.H{"message": "OK", "user": user, "token": token})
}
