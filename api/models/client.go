package models

import (
	"github.com/gorilla/websocket"

)

type Client struct {
    Conn *websocket.Conn
	Participant User
}