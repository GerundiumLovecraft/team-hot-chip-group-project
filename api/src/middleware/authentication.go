package middleware

import (
	"fmt"
	"net/http"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/auth"
	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")[7:]

	fmt.Println(tokenString)

	token, err := auth.DecodeToken(tokenString)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "auth error"})
		return
	}

	ctx.Set("userID", token.UserID)
	ctx.Next()
}
