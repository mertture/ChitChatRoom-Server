package controllers

import (
	"context"
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mertture/ChitChatRoom-Server/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (server *Server) CreateRoom(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	room := models.Room{}

	if err := c.BindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room.Prepare()
	err := room.Validate("create")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//hashing the password
	room.BeforeSave();

	count, err := server.DB.Collection("Room").CountDocuments(ctx, bson.M{"name": room.Name})
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking for the name"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This name already exists"})
		return
	}

	resultInsertionNumber, insertErr := server.DB.Collection("Room").InsertOne(ctx, room)
	if insertErr != nil {
		msg := fmt.Sprintf("Room was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, resultInsertionNumber)
}

func (server *Server) GetRoomByID(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	room := models.Room{}

	if err := c.BindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := room.ID;

    err := server.DB.Collection("Room").FindOne(ctx, bson.M{"_id": roomID}).Decode(&room)
    if err != nil {
         c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
         return
    }

    c.JSON(http.StatusOK, room)
}

func (server *Server) ListRooms(c *gin.Context) {
    var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
    defer cancel()

    // Find all rooms in the Room collection
    cursor, err := server.DB.Collection("Room").Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the rooms"})
        return
    }

    // Decode all rooms into a slice
    var rooms []models.Room
    if err := cursor.All(ctx, &rooms); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while decoding the rooms"})
        return
    }

    // Send the slice of rooms as the JSON response
    c.JSON(http.StatusOK, rooms)
}