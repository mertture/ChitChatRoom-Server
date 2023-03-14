package controllers

import (
	"github.com/mertture/ChitChatRoom-Server/api/middlewares"
)

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.GET("/api/dashboard", middlewares.SetMiddlewareJSON(s.Dashboard))

	s.Router.POST("/api/user/register", middlewares.SetMiddlewareJSON(s.Register))
	s.Router.POST("/api/user/login", middlewares.SetMiddlewareJSON(s.Login))

	s.Router.GET("/api/user/me", middlewares.SetMiddlewareAuthentication(s.GetUserByToken))
}
