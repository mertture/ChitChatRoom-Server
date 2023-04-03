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

type Message struct {
    Action       string `json:"action"`
    Room         string `json:"room"`
    Participant  string `json:"participant"`
    Data        interface{} `json:"data"`
}

type ChatMessage struct {
    Email        string `json:"email"`
    Content      string `json:"content"`
    CreatedAt    float64 `json:"createdAt"`
}



var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
        return true // allow any origin
    },
}

func (server *Server) addClientToRoom(Room string, client models.Client) {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	server.clients[Room][client.Conn] = client.Participant
}

func (server *Server) removeClientFromRoom(Room string, userObjectID primitive.ObjectID) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	for c, p := range server.clients[Room] {
        if (p.ID == userObjectID) {
            delete(server.clients[Room], c)
        }
    }
}

func (server *Server) getParticipantsArray(Room string) []models.User{
    server.mutex.Lock()
	defer server.mutex.Unlock()
    // Create a slice to store the users
    participantsMap := server.clients[Room]
    participants := make([]models.User, 0, len(participantsMap))

    // Iterate over the map and append each user to the slice
    participantIncludedMap := make(map[primitive.ObjectID]bool)
    for _, participant := range participantsMap {
        if _, ok := participantIncludedMap[participant.ID]; !ok {
            participants = append(participants, participant)
            participantIncludedMap[participant.ID] = true
        }
    }
    return participants
}


func (server *Server) websocketHandler(c *gin.Context) {
    var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
    defer cancel()

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
		fmt.Println("conn err:", err)
        c.AbortWithError(http.StatusInternalServerError, err)
		return
    }
    defer conn.Close()

    var message Message
    // handle incoming and outgoing WebSocket messages
    // using the conn object
    for {
        _, messageByte, err := conn.ReadMessage()
        json.Unmarshal(messageByte, &message)
        fmt.Println("mes:", message)

        Room := message.Room

        if (server.clients[Room] == nil) {
            server.clients[Room] = make(map[*websocket.Conn]models.User)
        }
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                fmt.Println("ws err:", err)
            }
            return
        }

        switch message.Action {
        
        case "join":

            ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
            defer cancel()
            // get the participant object from the user collection
            userObjectID, err := primitive.ObjectIDFromHex(message.Participant)
            if err != nil {
                fmt.Println("ws err: ", err)
                c.JSON(http.StatusNotFound, gin.H{"error": "Unable to convert Room to ObjectId"})
                return
            }

            filter := bson.M{"_id": userObjectID}
            var user models.User
            if err := server.DB.Collection("User").FindOne(ctx, filter).Decode(&user); err != nil {
                fmt.Println("error occurred: ", err)
                break
            }

            client := models.Client{
                Conn:        conn,
                Participant: user,
            }
            
            server.addClientToRoom(Room, client)


            participants := server.getParticipantsArray(Room)

            response := models.UsersResponse{
                Action: "users",
                Participants: participants,
            }

            // find the updated room document
            roomObjectId, err := primitive.ObjectIDFromHex(Room)
            if err != nil {
                fmt.Println("err on room object id", err)
            }
            // filter to find the room with the given id
            filter = bson.M{"_id": roomObjectId}
            var room models.Room
            err = server.DB.Collection("Room").FindOne(ctx, filter).Decode(&room)
            if err != nil {
                fmt.Println("err on getting room", err)
                break
            }

            // extract the messages field from the room document
            messages := room.Messages

            chatMessages := models.MessagesResponse{
                Action: "message",
                Messages: messages,
            }
            
    
            for c, k := range server.clients[Room] {
                if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                    fmt.Println("websocket connection closed")
                    break
                }
                if err := c.WriteJSON(response); err != nil {
                    fmt.Println("err on sending to clients: ", err);
                }
                if (k.ID == user.ID) {
                    fmt.Println("sending my messages to clients:")
                    if err := c.WriteJSON(chatMessages); err != nil {
                        fmt.Println("err on sending my messages to clients: ", err);
                    }
                }
            }
        
        case "leave":
            // get the participant object from the user collection
            userObjectID, err := primitive.ObjectIDFromHex(message.Participant)
            if err != nil {
                fmt.Println(err)
                c.JSON(http.StatusNotFound, gin.H{"error": "Unable to convert userID to ObjectId"})
                return
            }

            server.removeClientFromRoom(Room, userObjectID)
        
            participants := server.getParticipantsArray(Room)

            response := models.UsersResponse{
                Action: "users",
                Participants: participants,
            }

            for c := range server.clients[Room] {
                if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                    fmt.Println("websocket connection closed")
                    break
                }
                if err := c.WriteJSON(response); err != nil {
                    fmt.Println("err on sending to clients: ", err);
                    //break
                }
            }
        case "chat":
            // get the participant object from the user collection
            userObjectID, err := primitive.ObjectIDFromHex(message.Participant)
            if err != nil {
                fmt.Println(err)
                c.JSON(http.StatusNotFound, gin.H{"error": "Unable to convert userID to ObjectId"})
                return
            }

            var chatMessage ChatMessage
            if data, ok := message.Data.(map[string]interface{}); ok {
                // Use type assertion to convert the data to ChatMessage struct
                chatMessage = ChatMessage{
                    Email: data["email"].(string),
                    Content: data["content"].(string),
                    CreatedAt: data["createdAt"].(float64),
                }
            } else {
                // Handle error when type assertion fails
                fmt.Printf("Failed to convert message data to ChatMessage: %+v", message.Data)
                break;
            }
            fmt.Println("chat:", chatMessage)    
        
            //participants := server.getParticipantsArray(Room)
            // create a new message to add to the room
            sender := models.Sender {
                ID: userObjectID,
                Email: chatMessage.Email,
            }

            newMessage := models.Message{
                Sender: sender,
                Content: chatMessage.Content,
                CreatedAt: chatMessage.CreatedAt,
            }           
            newMessage.Prepare()
            
            
            roomObjectId, err := primitive.ObjectIDFromHex(Room)
            if err != nil {
                fmt.Println("err on room object id", err)
            }
            // filter to find the room with the given id
            filter := bson.M{"_id": roomObjectId}

            // update to add the new message to the messages array
            update := bson.M{
                "$push": bson.M{
                    "messages": newMessage,
                },
            }

            // update the room document in the database
            _, err = server.DB.Collection("Room").UpdateOne(ctx, filter, update)
            if err != nil {
                fmt.Println("err on update", err)
                break;
            }

            // find the updated room document
            var room models.Room
            err = server.DB.Collection("Room").FindOne(ctx, filter).Decode(&room)
            if err != nil {
                fmt.Println("err on getting room", err)
                break
            }

            // extract the messages field from the room document
            messages := room.Messages

            response := models.MessagesResponse{
                Action: "message",
                Messages: messages,
            }

            for c := range server.clients[Room] {
                if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                    fmt.Println("websocket connection closed")
                    break
                }
                if err := c.WriteJSON(response); err != nil {
                    fmt.Println("err on sending to clients: ", err);
                    //break
                }
            }
        }
            
    }
}




