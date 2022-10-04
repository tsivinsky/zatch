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
	Clicks    uint   `json:"clicks" gorm:"clicks,default:0"`
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
		err := db.Db.First(&url, "short_id", shortId).Error
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":      false,
				"message": "url not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":   true,
			"data": url,
		})
	})

	app.Delete("/api/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":      false,
				"message": "error parsing id",
			})
		}

		if id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":      false,
				"message": "id required after /api/",
			})
		}

		var url Url
		err = db.Db.Delete(&url, "id", id).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":      false,
				"message": "error deleting url",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":      true,
			"message": "deleted successfully",
		})
	})

	app.Get("/:shortId", func(c *fiber.Ctx) error {
		shortId := c.Params("shortId")
		if shortId == "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":      false,
				"message": "shortId parameter required",
			})
		}

		var url Url
		err := db.Db.First(&url, "short_id", shortId).Error
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":      false,
				"message": "url not found",
			})
		}

		go func() {
			url.Clicks = url.Clicks + 1
			db.Db.Save(&url)
		}()

		return c.Redirect(url.Url)
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
