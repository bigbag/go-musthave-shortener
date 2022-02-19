package url

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/bigbag/go-musthave-shortener/internal/utils"
)

type URLHandler struct {
	urlService URLService
	log        logrus.FieldLogger
}

func NewURLHandler(urlRoute fiber.Router, us URLService, l logrus.FieldLogger) {
	handler := &URLHandler{urlService: us, log: l}

	urlRoute.Post("/", handler.createURL)
	urlRoute.Get("/:shortID", handler.changeLocation)
}

func (h *URLHandler) createURL(c *fiber.Ctx) error {
	fullURL := string(c.Body()[:])
	if fullURL == "" {
		return utils.SendJSONError(c, fiber.StatusBadRequest, "Please specify a valid full url")
	}

	url, err := h.urlService.BuildURL(fullURL)
	if err != nil {
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	shortURL := fmt.Sprintf("%s/%s", c.BaseURL(), url.ShortID)
	return c.Status(fiber.StatusCreated).SendString(shortURL)
}

func (h *URLHandler) changeLocation(c *fiber.Ctx) error {
	shortID := c.Params("shortID")
	if shortID == "" {
		return utils.SendJSONError(c, fiber.StatusBadRequest, "Please specify a valid short id")
	}

	url, err := h.urlService.FetchURL(shortID)
	if err != nil {
		if url == nil {
			return utils.SendJSONError(c, fiber.StatusNotFound, err.Error())
		}
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	c.Location(url.FullURL)
	return c.Status(fiber.StatusTemporaryRedirect).SendString("")

}
