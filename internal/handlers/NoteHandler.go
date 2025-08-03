package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"net/http"
	url2 "net/url"
	"note1/internal/config"
	"note1/internal/models"
	"note1/internal/services"
	"strconv"
	"time"
)

type NoteHandler struct {
	ServiceContainer *services.ServicesContainer
}

func NewNoteHandler(serviceContainer *services.ServicesContainer) *NoteHandler {
	return &NoteHandler{ServiceContainer: serviceContainer}
}

// CreateNote создает новую заметку
// @Summary Создать заметку
// @Description Создает новую заметку для указанного пользователя
// @Tags notes
// @Accept json
// @Produce json
// @Param userID path string true "UUID пользователя"
// @Param note body models.Note true "Данные заметки"
// @Success 201 {object} map[string]interface{} "Note created"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /note/create/{userID} [post]
func CreateNoteHandler(h *NoteHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("userID")
		err := h.CreateNote(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}

func (h *NoteHandler) CreateNote(c *gin.Context, id string) error {

	idUuint, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	var note models.Note
	//распаковка запроса
	if err := c.ShouldBindJSON(&note); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	user, err := h.ServiceContainer.UserService.FindUserById(&idUuint)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return err
	}

	note.CreatedBy = user.UserID

	//создание уникального префикса для ноты
	note.FilePrefix = uuid.NewString()

	//Создание note в бд
	if err := h.ServiceContainer.NoteService.CreateNote(c, &note); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":     "Note created",
		"file_prefix": note.FilePrefix,
		"noteId":      note.ID,
	})

	// получаю юзера который создаёт (как получить юзера вопрос)
	// передаю бакет
	// отправляем файл в минайо

	//config.MinioClient.PutObject()

	return nil
}

// GetNoteByID получает заметку по ID
// @Summary Получить заметку по ID
// @Description Получает детали заметки по её идентификатору
// @Tags notes
// @Produce json
// @Param id path string true "ID заметки"
// @Success 200 {object} models.Note
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Note not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /note/{id} [get]
func GetNoteByIDHandler(h *NoteHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := h.GetNoteByID(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}

func (h *NoteHandler) GetNoteByID(c *gin.Context, id string) error {
	var note *models.Note

	noteID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	uintId := uint(noteID)

	note, err = h.ServiceContainer.NoteService.GetNoteByID(c, uintId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, note)
	return nil
}

// DeleteNoteById удаляет заметку по ID
// @Summary Удалить заметку
// @Description Удаляет заметку по её идентификатору
// @Tags notes
// @Produce json
// @Param id path string true "ID заметки"
// @Success 200 {object} string "Note deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /note/{id} [delete]
func DeleteNoteByIdHandler(h *NoteHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := h.DeleteNoteById(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}

func (h *NoteHandler) DeleteNoteById(c *gin.Context, id string) error {

	noteID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	uintID := uint(noteID)

	err = h.ServiceContainer.NoteService.DeleteNote(c, uintID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, "note deleted successfully")

	return nil
}

// AddFileToNote добавляет файлы к заметке
// @Summary Добавить файлы к заметке
// @Description Загружает файлы в MinIO и связывает их с заметкой
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param userID query string true "UUID пользователя"
// @Param filePrefix query string true "Префикс файлов заметки"
// @Param file formData file true "Файлы для загрузки"
// @Success 201 {object} string "Files added successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /note/files/upload [post]
func AddFileToNoteHandler(h *NoteHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h.AddFileToNote(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

func (h *NoteHandler) AddFileToNote(c *gin.Context) error {

	userIDStr := c.Query("userID")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no userID"})
		return errors.New("no userID")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
		return err
	}

	filePrefix := c.Query("filePrefix")
	if filePrefix == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no filePrefix"})
		return errors.New("no filePrefix")
	}

	//поиск юзера
	user, err := h.ServiceContainer.UserService.FindUserById(&userID)
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

// GetAllFilesForNote получает все файлы для заметки
// @Summary Получить все файлы заметки
// @Description Получает список файлов с пресигнутыми URL для указанной заметки
// @Tags files
// @Produce json
// @Param id path string true "ID заметки"
// @Success 200 {array} map[string]string "Список файлов с URL"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /note/files/get/{id} [get]
func GetAllFilesForNoteHandler(h *NoteHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt64, err := strconv.ParseUint(id, 10, 64)
		idUint := uint(idInt64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		err = h.GetAllFilesForNote(c, idUint)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}

func (h *NoteHandler) GetAllFilesForNote(c *gin.Context, noteID uint) error {
	note, err := h.ServiceContainer.NoteService.GetNoteByID(c, noteID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "note not found"})
		return err
	}

	user, err := h.ServiceContainer.UserService.FindUserById(&note.CreatedBy)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return err
	}

	objects := config.MinioClient.ListObjects(
		c,
		user.BucketName,
		minio.ListObjectsOptions{
			Prefix: note.FilePrefix,
		},
	)

	var filesUrl []map[string]string

	for object := range objects {
		if object.Err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": object.Err.Error()})
			return object.Err
		}
		url, err := config.MinioClient.PresignedGetObject(
			c,
			user.BucketName,
			object.Key,
			1*time.Hour,
			url2.Values{},
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return err
		}

		filesUrl = append(filesUrl, map[string]string{
			"filename": object.Key,
			"url":      url.String(),
		})
	}
	c.JSON(http.StatusOK, filesUrl)
	return nil

}

// GetNotesFrom получает список заметок с пагинацией
// @Summary Получить список заметок
// @Description Получает список всех заметок с пагинацией
// @Tags notes
// @Produce json
// @Param page query string false "Номер страницы" default(1)
// @Param limit query string false "Количество элементов на странице" default(10)
// @Success 200 {array} models.Note
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /note/list [get]
func GetNotesFromHandler(h *NoteHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h.GetNotesFrom(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}

func (h *NoteHandler) GetNotesFrom(c *gin.Context) error {

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(pageStr)
	if err != nil || pageInt < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return errors.New("invalid page number")
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil || limitInt < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid limit number"})
		return errors.New("invalid limit number")
	}

	notes, err := h.ServiceContainer.NoteService.GetNotesFrom(c, pageInt, limitInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	c.JSON(http.StatusOK, notes)
	return nil
}
