package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"note1/internal/models"
	"note1/internal/services"
)

type UserHandler struct {
	ServiceContainer *services.ServicesContainer
}

func NewUserHandler(serviceContainer *services.ServicesContainer) *UserHandler {
	return &UserHandler{ServiceContainer: serviceContainer}
}

// LoginRequest represents login request body
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Status string      `json:"status"`
	Body   interface{} `json:"body"`
}

// RegisterRequest represents register request body
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// RegisterResponse represents register response
type RegisterResponse struct {
	Status string      `json:"status"`
	Body   interface{} `json:"body"`
}

// LoginHandler godoc
// @Summary Авторизация пользователя
// @Description Выполняет авторизацию пользователя через email и пароль
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Данные для входа"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/login [post]
func LoginHandler(serviceContainer *services.ServicesContainer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Вызываем оригинальный метод
		userHandler := &UserHandler{ServiceContainer: serviceContainer}
		userHandler.Login(c, req.Email, req.Password)
	}
}

// RegisterHandler godoc
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя и создаёт бакет для файлов
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Данные пользователя"
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/register [post]
func RegisterHandler(serviceContainer *services.ServicesContainer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Вызываем оригинальный метод
		userHandler := &UserHandler{ServiceContainer: serviceContainer}
		userHandler.Register(c)
	}
}

// Login авторизация юзера
func (u *UserHandler) Login(c *gin.Context, email string, password string) {
	requestBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:7777/api/v1/login", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	for key, value := range resp.Header {
		for _, value := range value {
			c.Header(key, value)
		}
	}

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
}

// Register Регистрация юзера и создания бакета для файлов закреплённого за юзером
func (u *UserHandler) Register(c *gin.Context) error {

	var user models.Users
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	requestBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    user.Email,
		Password: user.Password,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Println(err)
		return err
	}
	req, err := http.NewRequest("POST", "http://localhost:7777/api/v1/register", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	c.JSON(200, gin.H{"status": resp.Status, "body": body})
	return nil
}
