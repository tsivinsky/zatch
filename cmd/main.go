package main

import (
	"log"
	"os"
	"strings"
	"url-shortener/internal/db"
	"url-shortener/internal/middleware"
	"url-shortener/internal/router"
	"url-shortener/internal/tasks"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = tasks.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer tasks.Asynq.Close()

	go func() {
		err = tasks.CreateServer()
		if err != nil {
			log.Fatal(err)
		}
	}()

	app := fiber.New()

	app.Get("/api/urls", router.GetAllUrls)
	app.Post("/api/shorten", router.CreateNewUrl)

	app.Get("/api/:shortId", middleware.ShouldHaveShortIdParam, router.GetUrlByShortId)
	app.Delete("/api/:id", router.DeleteUrlById)

	app.Get("/:shortId", middleware.ShouldHaveShortIdParam, router.FindUrlAndRedirect)

	port := getPort(":5000")
	log.Fatal(app.Listen(port))
}

func getPort(fallbackPort string) string {
	port := os.Getenv("PORT")

	if port == "" {
		port = fallbackPort
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return port
}
