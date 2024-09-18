package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewCorsMiddleware(origins []string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     ""
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	})
}
