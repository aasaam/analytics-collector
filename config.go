package main

import (
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type config struct {
	logger            *zerolog.Logger
	staticCacheTTL    uint
	clientSource      bool
	allowedMetricsIPs []net.IP
	collectorURL      *url.URL
}

func newConfig(
	logLevel string,
	staticCacheTTL uint,
	clientSource bool,
	collectorURL string,
	allowedMetricsIPs string,
) *config {

	c := config{
		staticCacheTTL: staticCacheTTL,
		clientSource:   clientSource,
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

	for _, ipString := range strings.Split(allowedMetricsIPs, ",") {
		ip := net.ParseIP(strings.TrimSpace(ipString))
		if ip != nil {
			c.allowedMetricsIPs = append(c.allowedMetricsIPs, ip)
		}
	}

	return &c
}

func (c *config) canAccessMetrics(ip net.IP) bool {
	for _, ipValid := range c.allowedMetricsIPs {
		if ipValid.Equal(ip) {
			return true
		}
	}
	return false
}

func (c *config) getLogger() *zerolog.Logger {
	return c.logger
}
