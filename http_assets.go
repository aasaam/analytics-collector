package main

import (
	_ "embed"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ineffassign: ignore
//go:embed embed/build/a.js
var embedADotJS []byte

// ineffassign: ignore
//go:embed embed/build/a.src.js
var embedASrcDotJS []byte

// ineffassign: ignore
//go:embed embed/build/l.js
var embedLDotJS []byte

// ineffassign: ignore
//go:embed embed/build/l.src.js
var embedLSrcDotJS []byte

// ineffassign: ignore
//go:embed embed/build/amp.json
var embedAmpDotJSON []byte

func httpAppAssets(
	app *fiber.App,
	conf *config,
) {
	staticCacheTTL := conf.staticCacheTTL

	embedADotJS = replaceCollectorURL(embedADotJS, conf.collectorURL)
	app.Get("/a.js", func(c *fiber.Ctx) error {
		if strings.Contains(c.Request().URI().String(), "?") {
			return httpErrorResponse(c, errorQueryStringDisabled)
		}
		staticCache(c, staticCacheTTL)
		c.Set(fiber.HeaderContentType, mimetypeJS)
		return c.Send(embedADotJS)
	})

	embedASrcDotJS = replaceCollectorURL(embedASrcDotJS, conf.collectorURL)
	app.Get("/a.src.js", func(c *fiber.Ctx) error {
		if strings.Contains(c.Request().URI().String(), "?") {
			return httpErrorResponse(c, errorQueryStringDisabled)
		}
		staticCache(c, staticCacheTTL)
		c.Set(fiber.HeaderContentType, mimetypeJS)
		return c.Send(embedASrcDotJS)
	})

	embedLDotJS = replaceCollectorURL(embedLDotJS, conf.collectorURL)
	app.Get("/l.js", func(c *fiber.Ctx) error {
		if strings.Contains(c.Request().URI().String(), "?") {
			return httpErrorResponse(c, errorQueryStringDisabled)
		}
		staticCache(c, staticCacheTTL)
		c.Set(fiber.HeaderContentType, mimetypeJS)
		return c.Send(embedLDotJS)
	})

	embedLSrcDotJS = replaceCollectorURL(embedLSrcDotJS, conf.collectorURL)
	app.Get("/l.src.js", func(c *fiber.Ctx) error {
		if strings.Contains(c.Request().URI().String(), "?") {
			return httpErrorResponse(c, errorQueryStringDisabled)
		}
		staticCache(c, staticCacheTTL)
		c.Set(fiber.HeaderContentType, mimetypeJS)
		return c.Send(embedLSrcDotJS)
	})

	embedAmpDotJSON = replaceCollectorURL(embedAmpDotJSON, conf.collectorURL)
	app.Get("/amp.json", func(c *fiber.Ctx) error {
		if strings.Contains(c.Request().URI().String(), "?") {
			return httpErrorResponse(c, errorQueryStringDisabled)
		}
		staticCache(c, staticCacheTTL)
		c.Set(fiber.HeaderContentType, mimetypeJS)
		return c.Send(embedAmpDotJSON)
	})
}
