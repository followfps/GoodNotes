package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"note1/internal/config"
	"note1/internal/models"
	"note1/internal/services"
	"os"
	"time"
)

type UserHandler struct {
	ServiceContaner *services.ServicesContainer
}

func NewUserHandler(serviceContainer *services.ServicesContainer) *UserHandler {
	return &UserHandler{ServiceContaner: serviceContainer}
}

// Login авторизация юзера
func (u *UserHandler) Login(c *gin.Context, email string, password string) bool {
	//Подгрузка соли из env
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	//поиск юзера по Email
	user, err := u.ServiceContaner.UserService.FindUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	//Создание и отправка JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
	return true
}

// Register Регистрация юзера и создания бакета для файлов закреплённого за юзером
func (u *UserHandler) Register(c *gin.Context) error {
	var user models.Users
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	// Проверка валидации email по виду (example@projectx.com)
	err = u.ServiceContaner.UserService.ValidationEmailCheck(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Incorrect Email": err.Error()})
		return err
	}

	// Проверка есть ли такой имейл в бд
	exists, err := u.ServiceContaner.UserService.EmailExists(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return err
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return err
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	user.Password = string(hashedPassword)

	//создание слуачйного имени бакета
	BucketName := uuid.NewString()

	//temp, err := config.MinioClient.BucketExists(c, BucketName)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return err
	//}
	//if temp {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Bucket exists"})
	//}

	if temp, _ := config.MinioClient.BucketExists(c, BucketName); temp {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bucket exists"})
	} else {
		err = config.MinioClient.MakeBucket(c, BucketName, minio.MakeBucketOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error, cant create bucket for user": err.Error()})
		} else {
			user.BucketName = BucketName
			fmt.Println("Create bucket ok")
		}
	}

	//Создаётся юзер ID
	user.UserID = uuid.New()

	//Запись юзера в бд
	err = u.ServiceContaner.UserService.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	userIDTemp := user.UserID.String()
	c.JSON(http.StatusOK, gin.H{"user": user})
	c.JSON(http.StatusOK, gin.H{"user": userIDTemp})
	return nil
}

//func (u *UserHandler) GetAllUsers(c *gin.Context) *[]models.Users {
//	return &[]models.Users{}
//}
