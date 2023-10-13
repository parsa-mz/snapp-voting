package services

import (
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"os"
)

var MinioClient *minio.Client

func init() {
	godotenv.Load("local.env")

	var err error
	MinioClient, err = minio.New(os.Getenv("MINIO_STORAGE_ENDPOINT"), &minio.Options{
		Creds: credentials.NewStaticV4(os.Getenv("MINIO_STORAGE_ACCESS"),
			os.Getenv("MINIO_STORAGE_SECRET"), ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
}
