package controllers

import (
	"net/http"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/models"
	"github.com/gin-gonic/gin"
)

type JSONLeaderboardEntry struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	SpotsCreated int    `json:"spots_created"`
}

func GetLeaderboardData(ctx *gin.Context) {
	leaderboard, err := models.FetchLeaderboard()
	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
}