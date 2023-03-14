package controllers

import (
	"net/http"
	"github.com/mertture/ChitChatRoom-Server/api/responses"
)

func (server *Server) Dashboard(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")
}