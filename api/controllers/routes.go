package controllers

import (
	"github.com/mertture/ChitChatRoom-Server/api/middlewares"
)

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.GET("/", middlewares.SetMiddlewareJSON(s.Dashboard))

	s.Router.POST("/register", middlewares.SetMiddlewareJSON(s.Register))
	s.Router.POST("/login", middlewares.SetMiddlewareJSON(s.Login))
}
