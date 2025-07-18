package config

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

var MinioClient *minio.Client

func InitMinio() {

	//Инициализация клиента
	client, err := minio.New(os.Getenv("endpoint"), &minio.Options{
		Creds: credentials.NewStaticV4(os.Getenv("accessKey"), os.Getenv("secretKey"), ""),
	})
	if err != nil {
		log.Fatalln("Ошибка инициализации клиента", err)
		return
	}

	_, err = client.ListBuckets(context.Background())
	if err != nil {
		log.Fatalln("Cannot connect to Minio", err)
	}

	MinioClient = client
	fmt.Println("minio connected")

}
