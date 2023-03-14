package controllers

import (
	"github.com/mertture/ChitChatRoom-Server/api/middlewares"
	"github.com/gin-gonic/gin"
)

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.GET("/", gin.WrapH(middlewares.SetMiddlewareJSON(s.Dashboard)))

}
