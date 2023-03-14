package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/mertture/ChitChatRoom-Server/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (server *Server) Register(c *gin.Context) {
    
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        user := models.User{}

        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }

		user.Prepare()
		err := user.Validate("login")
		//hashing the password
		user.BeforeSave();
	
		if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            
        }

        count, err := server.DB.Collection("User").CountDocuments(ctx, bson.M{"email": user.Email})
        defer cancel()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
            
        }

        if count > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
            
        }

        resultInsertionNumber, insertErr := server.DB.Collection("User").InsertOne(ctx, user)
        if insertErr != nil {
            msg := fmt.Sprintf("User item was not created")
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
        }
        defer cancel()

        c.JSON(http.StatusOK, resultInsertionNumber)
    }
