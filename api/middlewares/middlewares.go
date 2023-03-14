
package middlewares

import (
	"errors"
	"net/http"
	"github.com/mertture/ChitChatRoom-Server/api/auth"
	"github.com/mertture/ChitChatRoom-Server/api/responses"
	"github.com/gin-gonic/gin"
)

func SetMiddlewareJSON(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		next(c)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}