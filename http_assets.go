package main

import (
	_ "embed"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ineffassign: ignore
//
//go:embed embed/build/a.js
var embedADotJS []byte

// ineffassign: ignore
//
//go:embed embed/build/a.src.js
var embedASrcDotJS []byte

// ineffassign: ignore
//
//go:embed embed/build/l.js
var embedLDotJS []byte

// ineffassign: ignore
//
//go:embed embed/build/l.src.js
var embedLSrcDotJS []byte

// ineffassign: ignore
//
//go:embed embed/build/amp.json
var embedAmpDotJSON []byte

func httpAppAssetsHelper(c *fiber.Ctx) *errorMessage {
	if strings.Contains(c.Request().URI().String(), "?") {
		return &errorQueryStringDisabled
	}

	version := c.Params("version")
	date, dateErr := time.Parse("20060102", version)
	if dateErr != nil {
		return &errorAssetsVersionFailed
	}

	min := time.Now().AddDate(0, 0, -3)
	max := time.Now().AddDate(0, 0, 3)

	if date.Before(min) || date.After(max) {
		return &errorAssetsVersionFailed
	}

	return nil
}

func httpAppAssets(
	app *fiber.App,
	conf *config,
) {
	embedADotJS = replaceCollectorURL(embedADotJS, conf.collectorURL)
	app.Get("/_/:version/a.js", func(c *fiber.Ctx) error {
		e := httpAppAssetsHelper(c)
		if e != nil {
			return httpErrorResponse(c, *e)
		}
		staticCache(c)
		c.Set(fiber.HeaderContentType, mimetypeJS)
		return c.Send(embedADotJS)
	})

	embedASrcDotJS = replaceCollectorURL(embedASrcDotJS, conf.collectorURL)
	app.Get("/_/:version/a.src.js", func(c *fiber.Ctx) error {
		e := httpAppAssetsHelper(c)
		if e != nil {
			return httpErrorResponse(c, *e)
		}
		staticCache(c)
		c.Set(fiber.HeaderContentType, mimetypeJS)
		return c.Send(embedASrcDotJS)
	})

	embedLDotJS = replaceCollectorURL(embedLDotJS, conf.collectorURL)
	app.Get("/_/:version/l.js", func(c *fiber.Ctx) error {
		e := httpAppAssetsHelper(c)
		if e != nil {
			return httpErrorResponse(c, *e)
		}
		staticCache(c)
		c.Set(fiber.HeaderContentType, mimetypeJS)
		return c.Send(embedLDotJS)
	})

	embedLSrcDotJS = replaceCollectorURL(embedLSrcDotJS, conf.collectorURL)
	app.Get("/_/:version/l.src.js", func(c *fiber.Ctx) error {
		e := httpAppAssetsHelper(c)
		if e != nil {
			return httpErrorResponse(c, *e)
		}
		staticCache(c)
		c.Set(fiber.HeaderContentType, mimetypeJS)
		return c.Send(embedLSrcDotJS)
	})

	embedAmpDotJSON = replaceCollectorURL(embedAmpDotJSON, conf.collectorURL)
	app.Get("/amp.json", func(c *fiber.Ctx) error {
		if strings.Contains(c.Request().URI().String(), "?") {
			return httpErrorResponse(c, errorQueryStringDisabled)
		}
		staticCacheLimit(c, conf.staticCacheTTL)
		c.Set(fiber.HeaderContentType, mimetypeJSON)
		return c.Send(embedAmpDotJSON)
	})

	app.Get("/_/:version/amp.json", func(c *fiber.Ctx) error {
		e := httpAppAssetsHelper(c)
		if e != nil {
			return httpErrorResponse(c, *e)
		}
		staticCache(c)
		c.Set(fiber.HeaderContentType, mimetypeJSON)
		return c.Send(embedAmpDotJSON)
	})
}
