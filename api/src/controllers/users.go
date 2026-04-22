package controllers

import (
	"net/http"

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
	id := ctx.Param("id")

	user, err := models.FindUser(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := models.FindUserByUsername(username)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Username not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}