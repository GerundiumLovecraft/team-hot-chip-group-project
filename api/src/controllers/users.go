package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/auth"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(ctx *gin.Context) {
	var newUser models.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	newUser.Username = strings.ToLower(strings.TrimSpace(newUser.Username))
	newUser.Email = strings.ToLower(strings.TrimSpace(newUser.Email))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		SendInternalError(ctx, err)
		return
	}
	newUser.HashedPassword = string(hashedPassword)

	savedUser, err := newUser.Save()

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // if error code is 23505 (duplicate), return 409 otherwise 500
			ctx.JSON(http.StatusConflict, gin.H{"message": "Email or username already exists"})
			return
		}
		SendInternalError(ctx, err)
		return
	}

	//create token on successful signup
	savedUserId := strconv.FormatUint(uint64(savedUser.ID), 10)
	token, _ := auth.GenerateToken(savedUserId)

	ctx.JSON(http.StatusCreated, gin.H{"message": "OK", "token": token})
}

func GetUserById(ctx *gin.Context) {
	// Get the user ID from the token
	userId, exists := ctx.Get("userID")
	userIdStr, ok := userId.(string)

	// Get the user ID from the parameters
	searchedId := ctx.Param("id")

	// Find the user by ID
	user, err := models.FindUser(searchedId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User with this ID not found"})
		return
	}

	// Check if the user's ID matches to the searched ID
	if !exists || !ok || userIdStr != searchedId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorised access to user's profile"})
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

func UpdateUserAvatar(ctx *gin.Context) {

	userId, _ := ctx.Get("userID")
	userIdStr := userId.(string)

	var body struct {
		Avatar string `json:"avatar" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "An image URL is required"})
		return
	}

	_, err := models.UpdateUserAvatar(userIdStr, body.Avatar)
	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	token, _ := auth.GenerateToken(userIdStr)

	updatedUser, err := models.FindUser(userIdStr)
	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        userIdStr,
			"username":  updatedUser.Username,
			"email":     updatedUser.Email,
			"createdAt": updatedUser.CreatedAt,
			"avatar":    updatedUser.Avatar,
		},
	})
}