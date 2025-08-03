package Middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note1/internal/auth"
	"strings"
)

func Middlewares() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header found"})
			return
		}

		if strings.HasPrefix(authHeader, "Bearer ") {
			authHeader = authHeader[7:]
		}

		userID, err := auth.ValidateToken(authHeader)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Set("UserID", userID)
		c.Next()
	}
}
