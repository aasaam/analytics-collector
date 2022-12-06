package main

import (
	"net/url"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type config struct {
	logger         *zerolog.Logger
	staticCacheTTL uint
	collectorURL   *url.URL
}

func newConfig(
	logLevel string,
	staticCacheTTL uint,
	collectorURL string,
) *config {

	c := config{
		staticCacheTTL: staticCacheTTL,
		collectorURL:   getURL(collectorURL),
	}

	// logger config
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	logConfigLevel, errLogLevel := zerolog.ParseLevel(logLevel)
	if errLogLevel == nil {
		zerolog.SetGlobalLevel(logConfigLevel)
	}
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	c.logger = &logger

	return &c
}

func (c *config) getLogger() *zerolog.Logger {
	return c.logger
}
