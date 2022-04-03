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

// ineffassign: ignore
//go:embed embed/build/a.js
var embedADotJS []byte

// ineffassign: ignore
//go:embed embed/build/l.js
var embedLDotJS []byte

// ineffassign: ignore
//go:embed embed/build/amp.json
var embedAmpDotJSON []byte

func replaceCollectorURL(in []byte, collectorURL *url.URL) []byte {
	str := string(in)
	str = strings.ReplaceAll(str, collector_url_replacement, collectorURL.String())
	return []byte(str)
}

func httpErrorResponse(c *fiber.Ctx, message interface{}, code int) error {
	defer prometheusResponseErrors.WithLabelValues(strconv.Itoa(code)).Inc()
	c.Status(code)
	return c.JSON(message)
}

func staticCache(c *fiber.Ctx, staticCacheTTL uint) {
	ttlString := strconv.FormatUint(uint64(staticCacheTTL), 10)
	c.Set(fiber.HeaderCacheControl, "public, max-age="+ttlString+", stale-while-revalidate="+ttlString)
	c.Set(nginx_x_accel_expires, ttlString)
}

func noCache(c *fiber.Ctx) {
	c.Set(fiber.HeaderCacheControl, "no-cache, no-store, must-revalidate")
	c.Set(fiber.HeaderPragma, "no-cache")
	c.Set(fiber.HeaderExpires, "0")
	c.Set(nginx_x_accel_expires, "0")
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

	ipObject := net.ParseIP(ipString)

	if ipObject == nil {
		panic("ip address is not found")
	}

	c.Locals("ip", ipObject)

	return ipObject
}

func newHTTPServer(
	conf *config,
	geoParser *geoParser,
	projectsManager *projects,
	st *storage,
) *fiber.App {
	refererParser := newRefererParser()
	userAgentParser := newUserAgentParser()

	staticCacheTTL := conf.staticCacheTTL

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

			defer prometheusResponseErrors.WithLabelValues(strconv.Itoa(code)).Inc()

			ip := getClientIP(c)

			defer conf.getLogger().
				Error().
				Str("type", error_type_app).
				Str("error", err.Error()).
				Str("ip", ip.String()).
				Str("method", c.Method()).
				Str("url", c.Request().URI().String()).
				Int("status_code", code).
				Send()

			return httpErrorResponse(c, code, code)
		},
	})

	// recover
	app.Use(recover.New())

	// cors
	app.Use(cors.New())

	// init middleware
	app.Use(func(c *fiber.Ctx) error {
		ip := getClientIP(c)

		if c.Path() == metrics_path {
			if !conf.canAccessMetrics(ip) {
				defer conf.getLogger().
					Warn().
					Str("type", error_type_app).
					Str("ip", ip.String()).
					Str("method", c.Method()).
					Str("url", c.Context().URI().String()).
					Msg("forbidden")
				return httpErrorResponse(c, 403, 403)
			}
			return c.Next()
		} else {
			defer prometheusTotalRequests.Inc()
		}

		defer conf.getLogger().
			Trace().
			Str("ip", ip.String()).
			Str("method", c.Method()).
			Str("url", c.Context().URI().String()).
			Msg("access granted")

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
			st,
		)
	})

	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		c.Status(404)
		return c.JSON("not found")
	})

	embedADotJS = replaceCollectorURL(embedADotJS, conf.collectorURL)
	app.Get("/a.js", func(c *fiber.Ctx) error {
		staticCache(c, staticCacheTTL)
		c.Set(fiber.HeaderContentType, mimetype_js)
		return c.Send(embedADotJS)
	})

	embedLDotJS = replaceCollectorURL(embedLDotJS, conf.collectorURL)
	app.Get("/l.js", func(c *fiber.Ctx) error {
		staticCache(c, staticCacheTTL)
		c.Set(fiber.HeaderContentType, mimetype_js)
		return c.Send(embedLDotJS)
	})

	embedAmpDotJSON = replaceCollectorURL(embedAmpDotJSON, conf.collectorURL)
	app.Get("/amp.json", func(c *fiber.Ctx) error {
		staticCache(c, staticCacheTTL)
		c.Set(fiber.HeaderContentType, mimetype_js)
		return c.Send(embedAmpDotJSON)
	})

	handler := promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{})
	app.Get(metrics_path, adaptor.HTTPHandler(handler))

	// 404
	app.Use(func(c *fiber.Ctx) error {
		defer prometheusResponseErrors.WithLabelValues("404").Inc()
		defer conf.getLogger().
			Debug().
			Str("type", error_type_app).
			Str("ip", getClientIP(c).String()).
			Str("method", c.Method()).
			Str("url", c.Context().URI().String()).
			Msg("not found")
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.StatusNotFound)
	})

	return app
}
