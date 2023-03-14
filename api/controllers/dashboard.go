package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func (server *Server) Dashboard(c *gin.Context) {
	c.JSON(http.StatusOK, "Welcome To This Awesome API")
}