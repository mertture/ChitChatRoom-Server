package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mertture/ChitChatRoom-Server/api/auth"
	"github.com/mertture/ChitChatRoom-Server/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (server *Server) GetUserByToken(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	auth.TokenValid(c)
    

    user := models.User{}
	stringID := c.MustGet("user").(string)
	userID, err := primitive.ObjectIDFromHex(stringID)
    err = server.DB.Collection("User").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
    if err != nil {
         c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
         return
    }

    c.JSON(http.StatusOK, user)
}