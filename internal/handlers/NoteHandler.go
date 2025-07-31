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

// @Summary Создание новой заметки
// @Description Создаёт новую заметку для пользователя
// @Tags notes
// @Accept json
// @Produce json
// @Param userID path string true "ID пользователя"
// @Success 201 {object} map[string]interface{} "Заметка создана"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 404 {string} string "Пользователь не найден"
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
	if err := h.ServiceContainer.NoteService.CreateNote(&note); err != nil {
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

// @Summary Получение заметки по ID
// @Description Возвращает заметку по её уникальному идентификатору
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "ID заметки"
// @Success 200 {object} models.Note "Заметка найдена"
// @Failure 404 {string} string "Заметка не найдена"
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

	note, err = h.ServiceContainer.NoteService.GetNoteByID(uintId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, note)
	return nil
}

// @Summary Удаление заметки по ID
// @Description Удаляет заметку по её уникальному идентификатору
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "ID заметки"
// @Success 200 {string} string "Заметка удалена"
// @Failure 404 {string} string "Заметка не найдена"
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

	err = h.ServiceContainer.NoteService.DeleteNote(uintID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, "note deleted successfully")

	return nil
}

// @Summary Загрузка файла к заметке
// @Description Загружает файл и привязывает его к заметке
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param userID formData string true "ID пользователя"
// @Param file formData file true "Файл"
// @Param filePrefix formData string true "Префикс файла"
// @Success 201 {string} string "Файлы успешно загружены"
// @Failure 400 {string} string "Ошибка загрузки файла"
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

// @Summary Получение всех файлов к заметке
// @Description Возвращает все файлы, прикреплённые к заметке
// @Tags files
// @Accept json
// @Produce json
// @Param id path uint true "ID заметки"
// @Success 200 {array} map[string]string "Список файлов с URL"
// @Failure 404 {string} string "Файлы не найдены"
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
	note, err := h.ServiceContainer.NoteService.GetNoteByID(noteID)
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

// @Summary Список указанного количества заметок
// @Description Возвращает список всех заметок
// @Tags notes
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество записей на странице"
// @Success 200 {array} models.Note "Список заметок"
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

	notes, err := h.ServiceContainer.NoteService.GetNotesFrom(pageInt, limitInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	c.JSON(http.StatusOK, notes)
	return nil
}
