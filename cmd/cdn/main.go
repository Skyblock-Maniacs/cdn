package main

import (
	"os"

	"github.com/Skyblock-Maniacs/cdn/cdn"
	"github.com/Skyblock-Maniacs/cdn/internal/logger"
	"github.com/Skyblock-Maniacs/cdn/internal/s3"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Error.Fatal("Error loading .env file: ", err)
	}

	s3 := s3.Connect(s3.Config{
		Region: os.Getenv("AWS_REGION"),
		EndpointURL: os.Getenv("AWS_ENDPOINT_URL"),
		AccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	})

	cdn.Run(s3)
}
