package controllers

import (
	"context"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/mertture/ChitChatRoom-Server/api/auth"
	"github.com/mertture/ChitChatRoom-Server/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)



func (server *Server) Login(c *gin.Context) {
	user := models.User{}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user.Prepare()

	if err := user.Validate("login"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})	
	}        
	c.JSON(http.StatusOK, token)
}


func (server *Server) SignIn(email, password string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var err error

	user := models.User{}

	err = server.DB.Collection("Users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken((user.ID).String())
}
