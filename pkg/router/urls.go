package router

import (
	"fmt"
	"time"
	"url-shortener/pkg/db"
	"url-shortener/pkg/tasks"

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

	if id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"message": "id required after /api/",
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

		urlId = urlBody.Name
	}

	newUrl := db.Url{
		ShortId: urlId,
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

	return c.JSON(newUrl)
}
