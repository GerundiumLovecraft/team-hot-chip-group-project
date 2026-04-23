package controllers

import (
	"net/http"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/auth"
	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/models"
	"github.com/gin-gonic/gin"
)

type JSONPost struct {
	ID      uint   `json:"_id"`
	Message string `json:"message"`
}

func GetAllPosts(ctx *gin.Context) {
	posts, err := models.FetchAllPosts()

	if err != nil {
		SendInternalError(ctx, err)
		return
	}

		// Convert posts to JSON Structs
	jsonPosts := make([]JSONPost, 0)
	for _, post := range *posts {
		jsonPosts = append(jsonPosts, JSONPost{
			Message: post.Message,
			ID:      post.ID,
		})
	}

	val, exists := ctx.Get("userID")
	if exists {
		userID, ok := val.(string)
		if ok {
			token, _ := auth.GenerateToken(userID)
			ctx.JSON(http.StatusOK, gin.H{"posts": jsonPosts, "token": token})
			return
	}
}

	ctx.JSON(http.StatusOK, gin.H{"posts": jsonPosts})
}
type createPostRequestBody struct {
	Message string
}

func CreatePost(ctx *gin.Context) {
	var requestBody createPostRequestBody
	err := ctx.BindJSON(&requestBody)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	if len(requestBody.Message) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Post message empty"})
		return
	}

	newPost := models.Post{
		Message: requestBody.Message,
	}

	_, err = newPost.Save()
	if err != nil {
		SendInternalError(ctx, err)
		return
	}

	val, _ := ctx.Get("userID")
	userID := val.(string)
	token, _ := auth.GenerateToken(userID)

	ctx.JSON(http.StatusCreated, gin.H{"message": "Post created", "token": token})
}
