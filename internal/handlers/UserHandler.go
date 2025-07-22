package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
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

// Login авторизация юзера
func (u *UserHandler) Login(c *gin.Context, email string, password string) {
	////Подгрузка соли из env
	//jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	////поиск юзера по Email
	//user, err := u.ServiceContainer.UserService.FindUserByEmail(email)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return false
	//}
	//// Проверка пароля
	//err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return false
	//}
	////Создание и отправка JWT токена
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	"email": email,
	//	"exp":   time.Now().Add(time.Hour * 24).Unix(),
	//})
	//tokenString, err := token.SignedString([]byte(jwtSecretKey))
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//}
	//c.JSON(http.StatusOK, gin.H{"token": tokenString})
	//return true

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(200, gin.H{"status": resp.Status, "body": body})

}

// Register Регистрация юзера и создания бакета для файлов закреплённого за юзером
func (u *UserHandler) Register(c *gin.Context) error {
	//var user models.Users
	//err := c.ShouldBindJSON(&user)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return err
	//}
	//
	//// Проверка валидации email по виду (example@projectx.com)
	//err = u.ServiceContainer.UserService.ValidationEmailCheck(user.Email)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"Incorrect Email": err.Error()})
	//	return err
	//}
	//
	//// Проверка есть ли такой имейл в бд
	//exists, err := u.ServiceContainer.UserService.EmailExists(user.Email)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	//	return err
	//}
	//if exists {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
	//	return err
	//}
	//
	//// Хэширование пароля
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//}
	//user.Password = string(hashedPassword)
	//
	////создание слуачйного имени бакета
	//BucketName := uuid.NewString()
	//
	//if temp, _ := config.MinioClient.BucketExists(c, BucketName); temp {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Bucket exists"})
	//} else {
	//	err = config.MinioClient.MakeBucket(c, BucketName, minio.MakeBucketOptions{})
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error, cant create bucket for user": err.Error()})
	//	} else {
	//		user.BucketName = BucketName
	//		fmt.Println("Create bucket ok")
	//	}
	//}
	//
	////Создаётся юзер ID
	//user.UserID = uuid.New()
	//
	////Запись юзера в бд
	//err = u.ServiceContainer.UserService.CreateUser(&user)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return err
	//}
	//
	//userIDTemp := user.UserID.String()
	//c.JSON(http.StatusOK, gin.H{"user": user})
	//c.JSON(http.StatusOK, gin.H{"user": userIDTemp})
	//return nil

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

//func (u *UserHandler) GetAllUsers(c *gin.Context) *[]models.Users {
//	return &[]models.Users{}
//}
