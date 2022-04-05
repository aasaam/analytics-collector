package main

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

const sampleUserAgent = "Mozilla/5.0 (Linux; Android 5.0.2; SAMSUNG SM-A500FU Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/3.3 Chrome/38.0.2125.102 Mobile Safari/537.36"

func TestHTTP1(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

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
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

	statics := []string{
		"/a.js",
		"/l.js",
		"/amp.json",
		"/robots.txt",
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
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

	r404 := []string{
		"/favicon.ico",
		"/ensure-not-exist",
	}

	for _, p := range r404 {
		r := httptest.NewRequest("GET", p, nil)
		rs, _ := app.Test(r)

		if rs.StatusCode < 400 && rs.StatusCode >= 500 {
			t.Errorf("invalid response")
		}
	}

}
func TestHTTP3(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

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
func TestHTTP10(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

	r1 := httptest.NewRequest("PATCH", "/?m=pv_ins&i=000000000000", nil)
	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
	rs1, _ := app.Test(r1)

	if rs1.StatusCode == fiber.StatusOK {
		t.Errorf("invalid response")
	}
}
func TestHTTP12(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

	r1 := httptest.NewRequest("PATCH", "/?m=pv_ins&i=00000000000_", nil)
	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
	rs1, _ := app.Test(r1)

	if rs1.StatusCode == fiber.StatusOK {
		t.Errorf("invalid response")
	}
}
func TestHTTP20(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := getTestProjects()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

	r1 := httptest.NewRequest("GET", "/?m=pv_ins&i=000000000000&u=https%3A%2F%2Fexample.com%2F", nil)
	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
	rs1, _ := app.Test(r1)

	if rs1.StatusCode != fiber.StatusOK {
		t.Errorf("invalid response")
	}
}
func TestHTTP21(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := getTestProjects()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

	r1 := httptest.NewRequest("GET", "/?m=pv_ins&i=000000000000", nil)
	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
	r1.Header.Set(fiber.HeaderReferer, "http://example.com")
	rs1, _ := app.Test(r1)

	if rs1.StatusCode != fiber.StatusOK {
		t.Errorf("invalid response")
	}
}
func TestHTTP22(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := getTestProjects()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

	r1 := httptest.NewRequest("GET", "/?m=pv_ins&i=000000000000", nil)
	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
	rs1, _ := app.Test(r1)

	if rs1.StatusCode == fiber.StatusOK {
		t.Errorf("invalid response")
	}
}

func TestHTTP30(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := getTestProjects()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

	r1 := httptest.NewRequest("POST", "/?m=err&i=000000000000&u=https%3A%2F%2Fexample.com%2F", strings.NewReader(`{"foo":true"}`))
	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
	r1.Header.Set(fiber.HeaderContentType, "application/json")
	rs1, _ := app.Test(r1)

	if rs1.StatusCode == fiber.StatusOK {
		t.Errorf("invalid response")
	}
}

func TestHTTP31(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := getTestProjects()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

	postData := postRequest{
		ClientErrorMessage: "msg",
		ClientErrorObject:  "errObject",
	}

	b, _ := json.Marshal(postData)

	r1 := httptest.NewRequest("POST", "/?m=err&i=000000000000&u=https%3A%2F%2Fexample.com%2F", strings.NewReader(string(b)))
	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
	r1.Header.Set(fiber.HeaderContentType, "application/json")
	rs1, _ := app.Test(r1)

	if rs1.StatusCode != fiber.StatusOK {
		t.Errorf("invalid response")
	}
}
func TestHTTP32(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := getTestProjects()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

	postData := postRequest{}

	b, _ := json.Marshal(postData)

	r1 := httptest.NewRequest("POST", "/?m=e_api&i=000000000000", strings.NewReader(string(b)))
	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
	r1.Header.Set(fiber.HeaderContentType, "application/json")
	rs1, _ := app.Test(r1)

	if rs1.StatusCode == fiber.StatusOK {
		t.Errorf("invalid response")
	}
}
func TestHTTP33(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := getTestProjects()
	storage := newStorage()
	app := newHTTPServer(c1, geoParser, projectsManager, storage)

	api := postRequestAPI{
		PrivateInstanceKey: "000000000000111111111111",
		ClientIP:           "8.8.8.8",
		ClientUserAgent:    "curl 1.1.2",
		ClientTime:         time.Now().Unix(),
	}

	ev1 := postRequestEvent{
		Category: "cat",
		Action:   "act",
		Label:    "lab",
		Value:    1,
	}
	ev2 := postRequestEvent{
		Category: "!@#",
		Action:   "!@#",
		Label:    "lab",
		Value:    1,
	}

	events := []postRequestEvent{ev1, ev2}

	postData := postRequest{
		API:    &api,
		Events: &events,
	}

	b, _ := json.Marshal(postData)

	r1 := httptest.NewRequest("POST", "/?m=e_api&i=000000000000", strings.NewReader(string(b)))
	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
	r1.Header.Set(fiber.HeaderContentType, "application/json")
	rs1, _ := app.Test(r1)

	if rs1.StatusCode != fiber.StatusOK {
		t.Errorf("invalid response")
	}
}
