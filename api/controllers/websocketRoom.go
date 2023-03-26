package controllers

import (
	"fmt"
    "time"
    "context"
	"net/http"
    "encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/mertture/ChitChatRoom-Server/api/models"
)

type messagePayload struct {
    Message     string `json:"message" bson:"message"`
    Participant string `json:"participant" bson:"participant"`
};


var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
        return true // allow any origin
    },
}

func (server *Server) websocketHandler(c *gin.Context) {

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
		fmt.Println("conn err:", err)
        c.AbortWithError(http.StatusInternalServerError, err)
		return
    }
    defer conn.Close()

    roomID := c.Param("roomid");

    if (server.clients[roomID] == nil) {
        server.clients[roomID] = make([]models.Client, 0)
    }

    var messagePayload messagePayload
    // handle incoming and outgoing WebSocket messages
    // using the conn object
    for {
        _, messageByte, err := conn.ReadMessage()
        json.Unmarshal(messageByte, &messagePayload)


        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                fmt.Println("ws err:", err)
            }
            return
        }
        message := messagePayload.Message
        if (message == "new participant joined") {

            ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
            defer cancel()

            // get the participant object from the user collection
            userObjectID, err := primitive.ObjectIDFromHex(messagePayload.Participant)
            if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Unable to convert roomID to ObjectId"})
                return
           }

            filter := bson.M{"_id": userObjectID}
            var user models.User
            if err := server.DB.Collection("User").FindOne(ctx, filter).Decode(&user); err != nil {
                fmt.Println("error occurred: ", err)
                break
            }

            response := models.NewParticipantJoinedResponse{
                Message: message,
                Participant: user,
            };

            client := models.Client{
                Conn:        conn,
                Participant: user,
            }

            // Append the new client to the slice of clients for this WebSocket connection
            server.clients[roomID] = append(server.clients[roomID], client)
            fmt.Println("room pushed:", server.clients[roomID])

            for _, c := range server.clients[roomID] {
                if err := c.Conn.WriteJSON(response); err != nil {
                    fmt.Println("err on sending to clients: ", err);
                    return
                }
            }
        } else if (message == "participant left") {
            ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
            defer cancel()

            // get the participant object from the user collection
            userObjectID, err := primitive.ObjectIDFromHex(messagePayload.Participant)
            if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Unable to convert userID to ObjectId"})
                return
           }

            roomObjectID, err := primitive.ObjectIDFromHex(roomID)
            if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Unable to convert roomID to ObjectId"})
                return
           }

            filter := bson.M{"_id": roomObjectID}

            // Define the update to remove the participant with the given ID from the participants array
            update := bson.M{"$pull": bson.M{"participants": userObjectID}}

            // Update the room object in MongoDB
            result, err := server.DB.Collection("Room").UpdateOne(ctx, filter, update)
            if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Unable delete participant"})
                return
            }

            // Check the number of documents matched and modified to verify the update was successful
            if result.MatchedCount != 1 || result.ModifiedCount != 1 {
                // Handle the case where the room with the given ID or participant with the given ID was not found
                fmt.Printf("could not remove participant %s from room %s", userObjectID, roomID)
                return
            }

            response := models.ParticipantLeftResponse{
                Message: message,
                Participant: messagePayload.Participant,
            };

            for idx, c := range server.clients[roomID] {
                if (c.Participant.ID == userObjectID) {
                    server.clients[roomID] = append(server.clients[roomID][:idx], server.clients[roomID][idx+1:]...)
                }
            }

            for _, c := range server.clients[roomID] {
                if err := c.Conn.WriteJSON(response); err != nil {
                    fmt.Println("err on sending to clients: ", err);
                    break
                }
            }
        }
    }
	return
}




