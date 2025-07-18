package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note1/internal/models"
	"note1/internal/services"
	"strconv"
)

type NoteHandler struct {
	ServiceContaner *services.ServicesContainer
}

func NewNoteHandler(serviceContainer *services.ServicesContainer) *NoteHandler {
	return &NoteHandler{ServiceContaner: serviceContainer}
}

func (h *NoteHandler) CreateNote(c *gin.Context) error {
	var note models.Note
	//распаковка запроса
	if err := c.ShouldBindJSON(&note); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	//Создание note в бд
	if err := h.ServiceContaner.NoteService.CreateNote(&note); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	c.JSON(http.StatusCreated, "note created")

	// получаю юзера который создаёт (как получить юзера вопрос)
	// передаю бакет
	// отправляем файл в минайо

	//config.MinioClient.PutObject()

	return nil
}

func (h *NoteHandler) GetNoteByID(c *gin.Context, id string) error {
	var note *models.Note

	noteID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	uintId := uint(noteID)

	note, err = h.ServiceContaner.NoteService.GetNoteByID(uintId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, note)
	return nil
}

func (h *NoteHandler) DeleteNoteById(c *gin.Context, id string) error {

	noteID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	uintID := uint(noteID)

	err = h.ServiceContaner.NoteService.DeleteNote(uintID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, "note deleted successfully")

	return nil
}
