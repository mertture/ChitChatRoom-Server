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


	s.Router.POST("/api/room/create", middlewares.SetMiddlewareAuthentication(s.CreateRoom))
	s.Router.GET("/api/room/:roomid", middlewares.SetMiddlewareAuthentication(s.GetRoomByID))
	s.Router.POST("/api/room/:roomid", middlewares.SetMiddlewareAuthentication(s.EnterRoomByPassword))
	s.Router.GET("/api/rooms", middlewares.SetMiddlewareAuthentication(s.ListRooms))


}
