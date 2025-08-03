package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"note1/internal/models"
	"os"
)

var DBNote *gorm.DB
var DBUsers *gorm.DB
var DBTags *gorm.DB

func InitDB() {

	//путь до бд
	//dsn := "root:12345@tcp(127.0.0.1:3306)/note?charset=utf8mb4&parseTime=True&loc=Local"

	dsn := os.Getenv("DSN")

	var err error

	DBNote, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	fmt.Println("database note connect success")

	err = DBNote.AutoMigrate(&models.Note{})

	if err != nil {
		log.Println("failed to auto migrate database Note: " + err.Error())
	}
	fmt.Println("database Notes auto migrate success")

	dsn2 := "root:12345@tcp(127.0.0.1:3306)/users?charset=utf8mb4&parseTime=True&loc=Local"

	DBUsers, err = gorm.Open(mysql.Open(dsn2), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	fmt.Println("database users connect success")

	err = DBUsers.AutoMigrate(&models.Users{})

	if err != nil {
		log.Println("failed to auto migrate database Users: " + err.Error())
	}
	fmt.Println("database Users auto migrate success")

	DBTags, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	fmt.Println("database tags connect success")

	err = DBUsers.AutoMigrate(&models.Tags{})
	if err != nil {
		log.Println("failed to auto migrate database Tags: " + err.Error())
	}
	fmt.Println("database Tags auto migrate success")

	//RepoNote := repositories.NewGORMNoteRepository(DBNote)
	//RepoUser := repositories.NewGORMUserRepository(DBUsers)

}
