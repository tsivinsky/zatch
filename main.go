package main

import (
	"log"
	"url-shortener/pkg/db"
	"url-shortener/pkg/middleware"
	"url-shortener/pkg/router"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/api/urls", router.GetAllUrls)
	app.Post("/api/shorten", router.CreateNewUrl)

	app.Get("/api/:shortId", middleware.ShouldHaveShortIdParam, router.GetUrlByShortId)
	app.Delete("/api/:id", router.DeleteUrlById)

	app.Get("/:shortId", middleware.ShouldHaveShortIdParam, router.FindUrlAndRedirect)

	log.Fatal(app.Listen(":5000"))
}
