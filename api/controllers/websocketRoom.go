package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
        return true // allow any origin
    },
}

func (server *Server) websocketHandler(c *gin.Context) {
	fmt.Println("aaa")
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
		fmt.Println("conn err:", err)
        //c.AbortWithError(http.StatusInternalServerError, err)
		return
    }
    defer conn.Close()

    // handle incoming and outgoing WebSocket messages
    // using the conn object
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
			fmt.Println("ws err:", err)
            break
        }
		fmt.Println("ws message:", message)

		if err := conn.WriteJSON(message); err != nil {
			fmt.Println("error occurred: ", err)
			break
		}
    }
	return
}




