package Middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note1/internal/auth"
)

func Middlewares() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header found"})
		}

		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		//})
		//if err != nil {
		//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		//}
		//
		//if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//	c.Set("UserEmail", claims["email"])
		//} else {
		//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header found"})
		//	return
		//}

		userID, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Set("UserID", userID)
		c.Next()
	}
}
