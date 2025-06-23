package lib

import (
	"context"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
)

func UploadToCloudinary(file multipart.File, filename string) (string, error)  {
	errENV := godotenv.Load()
	if errENV != nil {
		log.Println("No .env file available")
	}

	cld, _ := cloudinary.NewFromParams(getEnv("CLOUDINARY_CLOUD", "your-cloud-name"), getEnv("CLOUDINARY_API_KEY", "your-api-key"), getEnv("CLOUDINARY_API_SECRET", "your-api-secret"))

	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
