package main

import (
	_ "embed"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func replaceCollectorURL(in []byte, collectorURL *url.URL) []byte {
	str := string(in)
	str = strings.ReplaceAll(str, collectorURLReplacement, collectorURL.String())
	return []byte(str)
}

func httpErrorResponse(c *fiber.Ctx, errMsg errorMessage) error {
	defer promMetricHTTPErrors.WithLabelValues(strconv.Itoa(errMsg.code)).Inc()
	c.Status(errMsg.code)
	return c.JSON(errMsg.msg)
}

func staticCache(c *fiber.Ctx, staticCacheTTL uint) {
	ttlString := strconv.FormatUint(uint64(staticCacheTTL), 10)
	c.Set(fiber.HeaderCacheControl, "public, max-age="+ttlString)
	c.Set(nginxXAccelExpires, ttlString)
}

func noCache(c *fiber.Ctx) {
	c.Set(fiber.HeaderCacheControl, "no-cache, no-store, must-revalidate")
	c.Set(fiber.HeaderPragma, "no-cache")
	c.Set(fiber.HeaderExpires, "0")
	c.Set(nginxXAccelExpires, "0")
}

func getClientIP(c *fiber.Ctx) net.IP {
	ipObjectStored := c.Locals("ip")
	if ipObjectStored != nil {
		return ipObjectStored.(net.IP)
	}

	ipString := ""

	if ipString == "" {
		ipString = c.IP()
		ipStrings := c.IPs()
		if len(ipStrings) > 0 {
			ipString = ipStrings[0]
		}
	}

	var ipObject net.IP
	ipObject = net.ParseIP(ipString)
	if ipObject == nil {
		ipObject = net.ParseIP("0.0.0.0")
	}

	c.Locals("ip", ipObject)

	return ipObject
}

func newHTTPServer(
	conf *config,
	geoParser *geoParser,
	projectsManager *projects,
	storage *storage,
) *fiber.App {
	refererParser := newRefererParser()
	userAgentParser := newUserAgentParser()
	promRegistry := getPrometheusRegistry()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		StrictRouting:         true,
		Prefork:               false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			defer promMetricHTTPErrors.WithLabelValues(strconv.Itoa(code)).Inc()

			ip := getClientIP(c)

			defer conf.getLogger().
				Error().
				Str("type", errorTypeApp).
				Str("error", err.Error()).
				Str("ip", ip.String()).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Int("status_code", code).
				Send()

			return httpErrorResponse(c, errorInternalServerError)
		},
	})

	// recover
	app.Use(recover.New())

	// cors
	app.Use(cors.New())

	// init middleware
	app.Use(func(c *fiber.Ctx) error {
		ip := getClientIP(c)

		if c.Path() == metricsPath {
			if !conf.canAccessMetrics(ip) {
				code := fiber.StatusForbidden
				defer conf.getLogger().
					Warn().
					Str("type", errorTypeApp).
					Str("ip", ip.String()).
					Str("method", c.Method()).
					Str("path", c.Path()).
					Int("status_code", code).
					Send()
				return httpErrorResponse(c, errorMetricsForbidden)
			}
			return c.Next()
		}

		defer promMetricHTTPTotalRequests.Inc()

		defer conf.getLogger().
			Trace().
			Str("ip", ip.String()).
			Str("method", c.Method()).
			Str("url", c.Context().URI().String()).
			Msg("http_access")

		return c.Next()
	})

	app.All("/", func(c *fiber.Ctx) error {
		return httpRecord(
			c,
			conf,
			refererParser,
			geoParser,
			userAgentParser,
			projectsManager,
			storage,
		)
	})

	httpAppAssets(app, conf)

	handler := promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{})
	app.Get(metricsPath, adaptor.HTTPHandler(handler))

	// 404
	app.Use(func(c *fiber.Ctx) error {
		code := fiber.StatusNotFound
		defer promMetricHTTPErrors.WithLabelValues(strconv.Itoa(code)).Inc()
		defer conf.getLogger().
			Debug().
			Str("type", errorTypeApp).
			Str("ip", getClientIP(c).String()).
			Str("method", c.Method()).
			Str("url", c.Context().URI().String()).
			Int("status_code", code).
			Send()
		c.Status(code)
		return c.JSON(code)
	})

	return app
}
