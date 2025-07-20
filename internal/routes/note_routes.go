package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
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

		noteGroup.Use(Middlewares.Middlewares()).POST("/create/:userID", func(c *gin.Context) {

			id := c.Param("userID")
			userID := uuid.MustParse(id)

			err := handler.CreateNote(c, userID)
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

		noteGroup.Use(Middlewares.Middlewares()).POST("/upload", func(c *gin.Context) {

			userIDStr := c.Query("userID")
			if userIDStr == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "no userID"})
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
			}

			filePrefix := c.Query("filePrefix")
			if filePrefix == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "no filePrefix"})
				return
			}

			err = handler.AddFileToNote(c, userID, filePrefix)
			if err != nil {
				return
			}

		})
	}
}
