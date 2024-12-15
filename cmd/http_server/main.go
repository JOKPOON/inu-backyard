package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/team-inu/inu-backyard/infrastructure/captcha"
	"github.com/team-inu/inu-backyard/infrastructure/database"
	"github.com/team-inu/inu-backyard/infrastructure/fiber"
	"github.com/team-inu/inu-backyard/internal/config"
	"github.com/team-inu/inu-backyard/internal/logger"
	"github.com/team-inu/inu-backyard/internal/utils/session"
)

func MustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return v
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Starting server...")
	fmt.Println(MustGetenv("SMTP_HOST"))

	var fiberConfig config.FiberServerConfig

	config.SetConfig(&fiberConfig)
	config.PrintConfig()

	zapLogger := logger.NewZapLogger()

	gormDB, err := database.NewGorm(&fiberConfig.Database)
	if err != nil {
		panic(err)
	}

	turnstile := captcha.NewTurnstile(fiberConfig.Client.Auth.Turnstile.SecretKey)

	session := session.NewSession()

	fiberServer := fiber.NewFiberServer(
		fiberConfig,
		gormDB,
		turnstile,
		zapLogger,
		session,
	)

	fiberServer.Run()
}
