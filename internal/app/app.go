package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"

	"github.com/bigbag/go-musthave-shortener/internal/config"
	"github.com/bigbag/go-musthave-shortener/internal/url"
	"github.com/bigbag/go-musthave-shortener/internal/utils"
)

type Server struct {
	l logrus.FieldLogger
	f *fiber.App
}

func New(l logrus.FieldLogger, cfg *config.Config) *Server {
	fiberCfg := fiber.Config{
		ReadTimeout: time.Second * cfg.Server.ReadTimeout,
		IdleTimeout: time.Second * cfg.Server.IdleTimeout,
		Immutable:   true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			l.WithError(err).Error("Unexpected API error")
			return utils.SendJSONError(ctx, fiber.StatusInternalServerError, err.Error())
		},
	}

	f := fiber.New(fiberCfg)

	f.Use(logger.New(logger.Config{
		Output: l.(*logrus.Logger).Writer(),
	}))

	urlRepository := url.NewURLRepository()
	urlService := url.NewURLService(urlRepository)
	url.NewURLHandler(f.Group(""), urlService, cfg, l)

	return &Server{l: l, f: f}
}

func (s *Server) Start(addr string) error {
	return s.f.Listen(addr)
}

func (s *Server) Stop() error {
	return s.f.Shutdown()
}
