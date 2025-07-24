package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"note1/internal/config"
	"note1/internal/routes"
	"note1/internal/services"
	"os"
)

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

	// Listen and serve on 0.0.0.0:8080
	err = http.ListenAndServe(os.Getenv("mainPort"), r)
	if err != nil {
		log.Fatal("Ошибка запуска сервера", err)
	}
}

// http protokol(metodi, parametri, kodi, ), rest api pochitat', obrabotka oshibok, mvc pochitat'

// papka services, git podkluchit, logirovanie

// minio cheto namutit
