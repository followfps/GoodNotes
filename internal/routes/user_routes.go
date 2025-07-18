package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note1/internal/handlers"
	"note1/internal/services"
)

func userRoutesSetup(r *gin.Engine, serviceContainer *services.ServicesContainer) {

	handler := handlers.NewUserHandler(serviceContainer)
	noteGroup := r.Group("/user")
	{
		noteGroup.POST("/register", func(c *gin.Context) {
			err := handler.Register(c)
			if err != nil {
				fmt.Println(err, "Cannot register user")
			}
		})

		noteGroup.POST("/login", func(c *gin.Context) {
			var input struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}

			err := c.ShouldBindJSON(&input)
			if err != nil {
				fmt.Println(err, "Cannot bind json")
			}

			handler.Login(c, input.Email, input.Password)
		})
	}

}
