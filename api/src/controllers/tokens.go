package controllers

import (
	"fmt"
	"net/http"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/auth"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type CreateTokenRequestBody struct {
	Email    string
	Password string
}

func CreateToken(ctx *gin.Context) {
	var input CreateTokenRequestBody
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	fmt.Println(input)

	user, err := models.FindUserByEmail(input.Email)
	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(input.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Password incorrect"})
		return
	}

	token, err := auth.GenerateToken(string(user.ID))
	if err != nil {
		SendInternalError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"token": token, "message": "OK"})
}
