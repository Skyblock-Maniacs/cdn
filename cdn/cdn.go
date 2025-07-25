package cdn

import (
	"os"

	"github.com/Skyblock-Maniacs/cdn/internal/logger"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
)

var app *fiber.App
var s3Client *s3.Client

func init() {
	logger.Info.Println("Starting CDN service...")
	app = fiber.New()
}

func Run(s3 *s3.Client) {
	s3Client = s3
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(fiber_logger.New())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/*", GetFileHandler)
	app.Post("/transcripts", middleware, PostTranscriptHandler)

	logger.Error.Fatal(app.Listen(":" + os.Getenv("PORT")))
}

func middleware(c *fiber.Ctx) error {
	if c.Get("Authorization") != os.Getenv("AUTH_TOKEN") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	return c.Next()
}
