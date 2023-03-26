package controllers

import (
	"context"
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mertture/ChitChatRoom-Server/api/models"
	"github.com/mertture/ChitChatRoom-Server/api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type RoomWithParticipants struct {
    ID        	   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Name      	   string             `json:"name,omitempty" bson:"name,omitempty"`
    Participants   []models.User      `json:"participants,omitempty" bson:"-"`
}

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

	roomID, err := primitive.ObjectIDFromHex(c.Param("roomid"));
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to convert roomID to ObjectId"})
		return
   }
	objectIDs := []primitive.ObjectID{}
    users := []models.User{}

	room := models.Room{}

    err = server.DB.Collection("Room").FindOne(ctx, bson.M{"_id": roomID}).Decode(&room)
    if err != nil {
         c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
         return
    }

	for _, participantID := range room.Participants {
        objectID, err := primitive.ObjectIDFromHex(participantID)
        if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
        }
        objectIDs = append(objectIDs, objectID)
    }

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	participants, err := server.DB.Collection("User").Find(ctx, filter);
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
    }
    defer participants.Close(ctx)

	for participants.Next(ctx) {
        var user models.User
        err := participants.Decode(&user)
        if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Error while decoding participant users"})
			return
        }
        users = append(users, user)
    }

    if err := participants.Err(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
    }

	roomWithParticipants := RoomWithParticipants{
		ID: room.ID,
		Name: room.Name,
		Participants: users,
	}
    c.JSON(http.StatusOK, roomWithParticipants)
}

func (server *Server) EnterRoomByPassword(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	roomRequestBody := models.Room{}

	if err := c.BindJSON(&roomRequestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID, err := primitive.ObjectIDFromHex(c.Param("roomid"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cannot convert roomID to objectId"})
		return
	}
	room := models.Room{}

    err = server.DB.Collection("Room").FindOne(ctx, bson.M{"_id": roomID}).Decode(&room)
    if err != nil {
         c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
         return
    }

	err = models.VerifyPassword(room.Password, roomRequestBody.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room password is wrong"})
		return 
	}

	stringID := c.MustGet("user").(string)
	userID, err := primitive.ObjectIDFromHex(stringID)

	
	if (!utils.Contains(room.Participants, stringID)) {
	// Define the update operation
		update := bson.M{
			"$push": bson.M{
				"participants": userID,
			},
		}

		// Execute the update operation
		roomResult, err := server.DB.Collection("Room").UpdateOne(ctx, bson.M{"_id": roomID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Room"})
			return
		}

		if roomResult.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		fmt.Println(roomResult);
	}	

	// Room updated successfully
	c.JSON(http.StatusOK, gin.H{"message": "Entered to the room successfully"})
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