package url

import (
	"github.com/bigbag/go-musthave-shortener/internal/config"
	"github.com/bigbag/go-musthave-shortener/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type URLHandler struct {
	urlService URLService
	log        logrus.FieldLogger
	cfg        *config.Config
}

func NewURLHandler(urlRoute fiber.Router, us URLService, cfg *config.Config, l logrus.FieldLogger) {
	handler := &URLHandler{urlService: us, log: l, cfg: cfg}

	urlRoute.Get("/ping", handler.getStatus)

	urlRoute.Post("/", handler.createShortURL)
	urlRoute.Post("/api/shorten", handler.createShortURLJson)
	urlRoute.Post("/api/shorten/batch", handler.createBatchOfShortURL)

	urlRoute.Get("/:shortID", handler.changeLocation)
	urlRoute.Get("/api/user/urls", handler.getUserURLs)
}

func (h *URLHandler) getBaseURL(c *fiber.Ctx) string {
	if h.cfg.BaseURL != "" {
		return h.cfg.BaseURL
	}
	return c.BaseURL()

}
func (h *URLHandler) getStatus(c *fiber.Ctx) error {
	err := h.urlService.Status()
	if err != nil {
		return utils.SendJSONError(
			c, fiber.StatusInternalServerError, "PG connection error",
		)
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"result": "OK"})
}

func (h *URLHandler) createShortURLJson(c *fiber.Ctx) error {
	req := new(JSONRequest)
	if err := c.BodyParser(req); err != nil {
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid full url",
		)
	}

	userID := c.Locals(h.cfg.UserContextKey).(string)
	url, err := h.urlService.BuildURL(h.getBaseURL(c), req.FullURL, userID)
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
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid full url",
		)
	}

	userID := c.Locals(h.cfg.UserContextKey).(string)
	url, err := h.urlService.BuildURL(h.getBaseURL(c), fullURL, userID)
	if err != nil {
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).SendString(url.ShortURL)
}

func (h *URLHandler) createBatchOfShortURL(c *fiber.Ctx) error {
	var items BatchRequest

	if err := c.BodyParser(&items); err != nil {
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid barch request",
		)
	}

	userID := c.Locals(h.cfg.UserContextKey).(string)
	result, err := h.urlService.BuildBatchOfURL(h.getBaseURL(c), items, userID)
	if err != nil {
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *URLHandler) changeLocation(c *fiber.Ctx) error {
	shortID := c.Params("shortID")
	if shortID == "" {
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid short id",
		)
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

func (h *URLHandler) getUserURLs(c *fiber.Ctx) error {
	userID := c.Locals(h.cfg.UserContextKey).(string)

	result, err := h.urlService.FetchUserURLs(h.getBaseURL(c), userID)
	if err != nil {
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	if len(result) == 0 {
		return utils.SendJSONError(c, fiber.StatusNoContent, "URLs not found")
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
