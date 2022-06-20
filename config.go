package main

import (
	"net/url"
	"os"

	"github.com/rs/zerolog"
)

type config struct {
	logger         *zerolog.Logger
	staticCacheTTL uint
	testMode       bool
	collectorURL   *url.URL
}

func newConfig(
	logLevel string,
	staticCacheTTL uint,
	testMode bool,
	collectorURL string,
) *config {

	c := config{
		staticCacheTTL: staticCacheTTL,
		testMode:       testMode,
		collectorURL:   getURL(collectorURL),
	}

	// logger config
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
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
