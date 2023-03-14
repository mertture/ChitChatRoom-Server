package api

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/mertture/ChitChatRoom-Server/api/controllers"
)

var server = controllers.Server{}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_URL"))

	server.Run(":8080")

}