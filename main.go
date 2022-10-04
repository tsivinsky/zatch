package main

import (
	"log"
	"time"
	"url-shortener/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/teris-io/shortid"
)

type Url struct {
	Id        uint   `json:"id" gorm:"primaryKey"`
	ShortId   string `json:"short_id" gorm:"short_id"`
	Url       string `json:"long_url" gorm:"url"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UrlBody struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

func main() {
	godotenv.Load()

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Db.AutoMigrate(&Url{})
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/api/urls", func(c *fiber.Ctx) error {
		var urls []Url
		db.Db.Find(&urls)

		return c.JSON(fiber.Map{
			"ok":   true,
			"data": urls,
		})
	})

	app.Get("/api/:shortId", func(c *fiber.Ctx) error {
		shortId := c.Params("shortId")
		if shortId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":      false,
				"message": "shortId required after /api/",
			})
		}

		var url Url
		db.Db.Where("short_id", shortId).Find(&url)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":   true,
			"data": url,
		})
	})

	app.Get("/:shortId", func(c *fiber.Ctx) error {
		shortId := c.Params("shortId")
		if shortId == "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":      false,
				"message": "url not found",
			})
		}

		var urls []Url
		db.Db.Find(&urls)

		for _, url := range urls {
			if url.ShortId == shortId {
				return c.Redirect(url.Url)
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":      false,
			"message": "url not found",
		})
	})

	app.Post("/api/shorten", func(c *fiber.Ctx) error {
		urlBody := new(UrlBody)

		if err := c.BodyParser(urlBody); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":      false,
				"message": err.Error(),
			})
		}

		urlId, err := shortid.Generate()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":      false,
				"message": err.Error(),
			})
		}

		var urls []Url
		db.Db.Find(&urls)

		if urlBody.Name != "" {
			for _, url := range urls {
				if url.ShortId == urlBody.Name {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"ok":      false,
						"message": "this url name already taken",
					})
				}
			}

			urlId = urlBody.Name
		}

		newUrl := Url{
			ShortId: urlId,
			Url:     urlBody.Url,
		}
		db.Db.Create(&newUrl)

		return c.JSON(newUrl)
	})

	log.Fatal(app.Listen(":5000"))
}
