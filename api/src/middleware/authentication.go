package middleware

import (
	"fmt"
	"net/http"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/auth"
	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware(ctx *gin.Context) {
    authHeader := ctx.GetHeader("Authorization")
    if len(authHeader) < 7 {
        ctx.JSON(http.StatusUnauthorized, gin.H{"message": "auth error"})
        ctx.Abort()
        return
    }
    tokenString := authHeader[7:]
    token, err := auth.DecodeToken(tokenString)
    if err != nil {
        fmt.Println(err)
        ctx.JSON(http.StatusUnauthorized, gin.H{"message": "auth error"})
        ctx.Abort()
        return
    }
    ctx.Set("userID", token.UserID)
    ctx.Next()
}
