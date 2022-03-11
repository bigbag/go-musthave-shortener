package url

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/bigbag/go-musthave-shortener/internal/config"
	"github.com/bigbag/go-musthave-shortener/internal/utils"
)

type URLHandler struct {
	urlService URLService
	log        logrus.FieldLogger
	cfg        *config.Config
}

func NewURLHandler(urlRoute fiber.Router, us URLService, cfg *config.Config, l logrus.FieldLogger) {
	handler := &URLHandler{urlService: us, log: l, cfg: cfg}

	urlRoute.Post("/", handler.createShortURL)
	urlRoute.Post("/api/shorten", handler.createShortURLJson)
	urlRoute.Get("/:shortID", handler.changeLocation)

}

func (h *URLHandler) getBaseURL(c *fiber.Ctx) string {
	if h.cfg.BaseURL != "" {
		return h.cfg.BaseURL
	}
	return c.BaseURL()

}

func (h *URLHandler) createShortURLJson(c *fiber.Ctx) error {
	req := new(URL)
	if err := c.BodyParser(req); err != nil {
		return utils.SendJSONError(c, fiber.StatusBadRequest, "Please specify a valid full url")
	}

	url, err := h.urlService.BuildURL(h.getBaseURL(c), req.FullURL)
	if err != nil {
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"result": url.ShortURL,
	})
}

func (h *URLHandler) createShortURL(c *fiber.Ctx) error {
	fullURL := string(c.Body())
	if fullURL == "" {
		return utils.SendJSONError(c, fiber.StatusBadRequest, "Please specify a valid full url")
	}

	url, err := h.urlService.BuildURL(h.getBaseURL(c), fullURL)
	if err != nil {
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).SendString(url.ShortURL)
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
