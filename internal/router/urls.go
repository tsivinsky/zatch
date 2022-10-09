package router

import (
	"fmt"
	"time"
	"url-shortener/internal/db"
	"url-shortener/internal/tasks"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/teris-io/shortid"
)

type UrlBody struct {
	Url        string `json:"url"`
	Name       string `json:"name"`
	AutoDelete uint   `json:"auto_delete"`
}

func GetAllUrls(c *fiber.Ctx) error {
	var urls []db.Url
	db.Db.Find(&urls)

	return c.JSON(fiber.Map{
		"ok":   true,
		"data": urls,
	})
}

func GetUrlByShortId(c *fiber.Ctx) error {
	shortId := c.Params("shortId")

	var url db.Url
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
}

func DeleteUrlById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": "error parsing id",
		})
	}

	var url db.Url
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
}

func FindUrlAndRedirect(c *fiber.Ctx) error {
	shortId := c.Params("shortId")

	var url db.Url
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
}

func CreateNewUrl(c *fiber.Ctx) error {
	var err error

	urlBody := new(UrlBody)

	if err := c.BodyParser(urlBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": err.Error(),
		})
	}

	urlShortId := ""

	var urls []db.Url
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

		urlShortId = urlBody.Name
	} else {
		urlShortId, err = shortid.Generate()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":      false,
				"message": err.Error(),
			})
		}
	}

	newUrl := db.Url{
		ShortId: urlShortId,
		Url:     urlBody.Url,
	}
	db.Db.Create(&newUrl)

	if urlBody.AutoDelete > 0 {
		task, err := tasks.NewAutoDeleteUrlTask(newUrl.Id)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
		_, err = tasks.Asynq.Enqueue(task, asynq.ProcessIn(time.Duration(urlBody.AutoDelete*uint(time.Minute))))
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
	}

	return c.Status(fiber.StatusCreated).JSON(newUrl)
}
