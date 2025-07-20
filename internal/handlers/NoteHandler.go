package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"net/http"
	"note1/internal/config"
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

func (h *NoteHandler) CreateNote(c *gin.Context, id uuid.UUID) error {
	var note models.Note
	//распаковка запроса
	if err := c.ShouldBindJSON(&note); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	user, err := h.ServiceContaner.UserService.FindUserById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return err
	}

	note.CreatedBy = user.UserID

	//создание уникального префикса для ноты
	note.FilePrefix = uuid.NewString()

	//Создание note в бд
	if err := h.ServiceContaner.NoteService.CreateNote(&note); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	c.JSON(http.StatusCreated, note.FilePrefix+"note created")

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

func (h *NoteHandler) AddFileToNote(c *gin.Context, userID uuid.UUID, filePrefix string) error {

	//поиск юзера
	user, err := h.ServiceContaner.UserService.FindUserById(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error, user not found": err.Error()})
		return err
	}

	// Проверка и парсинг multipart/form-data
	err = c.Request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form"})
		return err
	}

	//получиение файлов из тела запроса
	files := c.Request.MultipartForm.File["file"]

	//цикл для записи файлов в minio
	for _, fileHeader := range files {
		//открытие текущего файла
		file, err := fileHeader.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return err
		}
		//закрытие файла после загрузки
		defer file.Close()

		//помещение файла в minio
		_, err = config.MinioClient.PutObject(
			c,
			user.BucketName,
			fmt.Sprintf("%s_%s", filePrefix, fileHeader.Filename),
			file,
			fileHeader.Size,
			minio.PutObjectOptions{},
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return err
		}
	}
	c.AbortWithStatusJSON(http.StatusCreated, "files added successfully")
	return nil

}

//func (h *NoteHandler) GetAllFilesForNote(c *gin.Context, userID *uuid.UUID, filePrefix string) error {
//
//}
