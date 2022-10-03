package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/teris-io/shortid"
)

type Url struct {
	Id      int    `json:"id"`
	ShortId string `json:"short_id"`
	LongUrl string `json:"long_url"`
}

type UrlBody struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

var urls []Url = []Url{}

func main() {
	app := fiber.New()

	app.Get("/api/urls", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ok":   true,
			"data": urls,
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

		for _, url := range urls {
			if url.ShortId == shortId {
				return c.Redirect(url.LongUrl)
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
			Id:      len(urls),
			ShortId: urlId,
			LongUrl: urlBody.Url,
		}

		urls = append(urls, newUrl)

		return c.JSON(newUrl)
	})

	log.Fatal(app.Listen(":5000"))
}
