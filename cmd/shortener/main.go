package main

import (
	"io/ioutil"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/bigbag/go-musthave-shortener/internal/app"
	"github.com/bigbag/go-musthave-shortener/internal/config"
)

// global variables
var (
	cfg *config.Config
	log logrus.StdLogger
)

func initLogger(cfg *config.Config) {
	logger := logrus.New()

	if level, err := logrus.ParseLevel(cfg.Logger.Level); err == nil {
		log.Printf("Logger level set to %s\n", level.String())
		logger.SetLevel(level)
	} else {
		log.Printf("Failed to parse log level: %s\n", err.Error())
		logger.SetLevel(logrus.ErrorLevel)
	}

	switch cfg.Logger.Output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		logger.SetOutput(ioutil.Discard)
	}

	switch cfg.Logger.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	log = logger
}

func _init() {
	log = stdLog.New(os.Stdout, "", 0)

	var err error
	cfg, err = config.New()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v\n", err)
	}

	initLogger(cfg)
}

// shutdown implements all graceful shutdown logic
func shutdown(_ os.Signal, apiServer *app.Server) {
	log.Println("Shutting down...")
	if err := apiServer.Stop(); err != nil {
		log.Printf("Error stopping api: %v", err)
	}
	log.Println("Running cleanup tasks...")
	log.Println("Fiber was successful shutdown.")
	os.Exit(0)
}

func main() {
	_init()

	var (
		logger    = log.(logrus.FieldLogger)
		apiServer = app.New(logger, cfg)
	)

	// start HTTP API server
	go func() {
		log.Printf("API server listening on %s", cfg.Server.Listen)
		if err := apiServer.Start(cfg.Server.Listen); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start api server: %v", err)
		}
	}()

	// listen for exit signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	shutdown(<-c, apiServer)
}
