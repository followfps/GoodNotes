package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	_ "note1/docs"
	"note1/internal/config"
	"note1/internal/routes"
	"note1/internal/services"
	"os"
)

// @title Note App API
// @version 1.0
// @description API server for creating and reading notes
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {

	// Пока что env не фул надо заполнить
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	fmt.Println("env loaded")

	config.InitDB()
	config.InitMinio()

	r := gin.Default()

	routes.SetupRoutes(r, services.NewServicesContainer())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Listen and serve on 0.0.0.0:8080
	err = http.ListenAndServe(os.Getenv("mainPort"), r)
	if err != nil {
		log.Fatal("Ошибка запуска сервера", err)
	}
}

// http protokol(metodi, parametri, kodi, ), rest api pochitat', obrabotka oshibok, mvc pochitat'

// papka services, git podkluchit, logirovanie

// minio cheto namutit
