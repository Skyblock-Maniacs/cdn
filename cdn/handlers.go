package cdn

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetFileHandler(c *fiber.Ctx) error {
	ctx := context.Background()
	key := c.Params("*")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File key is required",
		})
	}

	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(c.Params("*")),
	})

	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey:") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "File not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve file",
			"message": err.Error(),
		})
	}
	defer result.Body.Close()
	

	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)

	return c.Type(*result.ContentType).Send(buf.Bytes())
}

func PostTranscriptHandler(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form data",
			"message": err.Error(),
		})
	}
	files, ok := form.File["file"]
	if !ok || len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}
	file := files[0]

	if !(strings.HasPrefix(file.Header.Get("content-type"), "text/html")) {
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
			"error": "Unsupported file type",
		})
	}

	extractedFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}
	defer extractedFile.Close()

	id := strings.Split(uuid.New().String(), "-")[0]
	
	extractedFileBytes, err := io.ReadAll(extractedFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read file",
			"message": err.Error(),
		})
	}

	ctx := context.Background()
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key: 	aws.String("transcripts/" + id + ".html"),
		ContentType: aws.String("text/html"),
		Body:   bytes.NewReader(extractedFileBytes),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"url": os.Getenv("URI") + "/transcripts/" + id + ".html",
	})
}
