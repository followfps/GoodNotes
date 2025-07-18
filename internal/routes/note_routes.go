package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note1/internal/Middlewares"
	"note1/internal/handlers"
	"note1/internal/services"
)

func noteRoutesSetup(r *gin.Engine, serviceContainer *services.ServicesContainer) {

	//repo := repositories.NewGORMNoteRepository(config.DBNote)
	//handler := handlers.NewNoteHandler(repo)

	handler := handlers.NewNoteHandler(serviceContainer)

	noteGroup := r.Group("/note")
	{
		noteGroup.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			err := handler.GetNoteByID(c, id)
			if err != nil {
				fmt.Println(err, "Note not found")
				return
			}

		})

		noteGroup.Use(Middlewares.Middlewares()).POST("/create", func(c *gin.Context) {
			err := handler.CreateNote(c)
			if err != nil {
				fmt.Println(err, "Note not found")
				return
			}

		})

		noteGroup.Use(Middlewares.Middlewares()).DELETE("/:id", func(c *gin.Context) {
			id := c.Param("id")
			err := handler.DeleteNoteById(c, id)
			if err != nil {
				fmt.Println(err, "Note not found")
				return
			}
		})
	}
}
