package api

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/mertture/ChitChatRoom-Server/api/controllers"
)

var server = controllers.Server{}

func Run() {

	if err:= godotenv.Load(); err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	}

	server.Initialize(os.Getenv("DB_URL"))

	server.Run(":8080")

}