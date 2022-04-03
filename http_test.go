package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHTTP1(t *testing.T) {
	c1 := newConfig("error", 0, true, "http://127.0.0.1", "")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	st := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, st)

	statics := []string{
		"/a.js",
		"/l.js",
	}

	for _, p := range statics {
		r := httptest.NewRequest("GET", p, nil)
		rs, _ := app.Test(r)

		if rs.StatusCode != fiber.StatusOK {
			t.Errorf("invalid response")
		}
	}
}
func TestHTTP11(t *testing.T) {
	c1 := newConfig("error", 0, true, "http://127.0.0.1", "")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	st := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, st)

	statics := []string{
		"/a.js",
		"/l.js",
		"/amp.json",
	}

	for _, p := range statics {
		r := httptest.NewRequest("GET", p, nil)
		rs, _ := app.Test(r)

		if rs.StatusCode != fiber.StatusOK {
			t.Errorf("invalid response")
		}
	}
}
func TestHTTP2(t *testing.T) {
	c1 := newConfig("error", 0, true, "http://127.0.0.1", "")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	st := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, st)

	r404 := []string{
		"/favicon.ico",
		"/ensure-not-exist",
	}

	for _, p := range r404 {
		r := httptest.NewRequest("GET", p, nil)
		rs, _ := app.Test(r)

		if rs.StatusCode != fiber.StatusNotFound {
			t.Errorf("invalid response")
		}
	}

}
func TestHTTP3(t *testing.T) {
	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	st := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, st)

	r1 := httptest.NewRequest("GET", "/metrics", nil)
	r1.Header.Set(fiber.HeaderXForwardedFor, "192.168.1.1")
	rs1, _ := app.Test(r1)

	if rs1.StatusCode == fiber.StatusOK {
		t.Errorf("invalid response")
	}

	r2 := httptest.NewRequest("GET", "/metrics", nil)
	r2.Header.Set(fiber.HeaderXForwardedFor, "127.0.0.1")
	rs2, _ := app.Test(r2)

	if rs2.StatusCode != fiber.StatusOK {
		t.Errorf("invalid response")
	}
}
