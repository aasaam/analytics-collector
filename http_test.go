package main

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	fake "github.com/brianvoe/gofakeit/v6"
)

func TestHTTPAssets(t *testing.T) {
	conf := newConfig("error", 0, "http://127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	refererParser := newRefererParser()
	userAgentParser := newUserAgentParser()
	redisClient, redisClientErr := redisGetClient("redis://127.0.0.1:6379/0")
	if redisClientErr != nil {
		t.Error(redisClientErr)
		return
	}

	app := newHTTPServer(
		conf,
		geoParser,
		refererParser,
		userAgentParser,
		projectsManager,
		redisClient,
	)

	rq404 := httptest.NewRequest("GET", "/ensure-not-exist-path", nil)
	rq404.Header.Set("X-Real-IP", "9.9.9.9")
	rs404, _ := app.Test(rq404)

	if rs404.StatusCode != fiber.StatusNotFound {
		t.Errorf("invalid response")
	}

	version := time.Now().Format("20060102")

	statics := []string{
		"/_/" + version + "/a.js",
		"/_/" + version + "/a.src.js",
		"/_/" + version + "/l.js",
		"/_/" + version + "/l.src.js",
		"/_/" + version + "/amp.json",
		"/amp.json",
	}

	for _, p := range statics {
		rq := httptest.NewRequest("GET", p, nil)
		rs, _ := app.Test(rq)

		if rs.StatusCode != fiber.StatusOK {
			t.Errorf("invalid response")
		}
	}

	staticsInvalidVersion := []string{
		"/_/A!/a.js",
		"/_/20010101/a.js",
		"/_/22010101/a.js",
		"/_/" + version + "/a.js?a=1",
		"/_/" + version + "/a.src.js?b=2",
		"/_/" + version + "/l.js?b=2",
		"/_/" + version + "/l.src.js?b=2",
		"/_/" + version + "/amp.json?b=2",
		"/amp.json?v=1",
	}

	for _, p := range staticsInvalidVersion {
		rq := httptest.NewRequest("GET", p, nil)
		rs, _ := app.Test(rq)

		if rs.StatusCode != fiber.StatusGone {
			t.Errorf("invalid response")
		}
	}

}
func TestHTTP003(t *testing.T) {
	conf := newConfig("error", 0, "http://127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := getTestProjects()
	refererParser := newRefererParser()
	userAgentParser := newUserAgentParser()
	redisClient, redisClientErr := redisGetClient("redis://127.0.0.1:6379/0")
	if redisClientErr != nil {
		t.Error(redisClientErr)
		return
	}

	app := newHTTPServer(
		conf,
		geoParser,
		refererParser,
		userAgentParser,
		projectsManager,
		redisClient,
	)

	rq0 := httptest.NewRequest("GET", "/?i=000000000001&m=ERROR&err=Sample", nil)
	rq0.Header.Set("X-Real-IP", fake.IPv4Address())
	rs0, _ := app.Test(rq0)

	if rs0.StatusCode != errorInvalidModeOrProjectPublicID.code {
		t.Errorf("invalid response")
	}

	rq1 := httptest.NewRequest("GET", "/?i=12&m=err_l&err=Sample", nil)
	rq1.Header.Set("X-Real-IP", fake.IPv4Address())
	rs1, _ := app.Test(rq1)

	if rs1.StatusCode != errorInvalidModeOrProjectPublicID.code {
		t.Errorf("invalid response")
	}
}

func TestHTTP000(t *testing.T) {
	conf := newConfig("error", 0, "http://127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := getTestProjects()
	refererParser := newRefererParser()
	userAgentParser := newUserAgentParser()
	redisClient, redisClientErr := redisGetClient("redis://127.0.0.1:6379/0")
	if redisClientErr != nil {
		t.Error(redisClientErr)
		return
	}

	app := newHTTPServer(
		conf,
		geoParser,
		refererParser,
		userAgentParser,
		projectsManager,
		redisClient,
	)

	rq0 := httptest.NewRequest("GET", "/?i=000000000000&m=err_l&err=Sample&u=http%3A//example-not.net/", nil)
	rq0.Header.Set("X-Real-IP", fake.IPv4Address())
	rs0, _ := app.Test(rq0)

	if rs0.StatusCode != fiber.StatusForbidden {
		t.Errorf("invalid response")
	}

	rq1 := httptest.NewRequest("GET", "/?i=000000000000&m=err_l&err=Sample&u=http%3A//example.net/", nil)
	rq1.Header.Set("X-Real-IP", fake.IPv4Address())
	rs1, _ := app.Test(rq1)

	if rs1.StatusCode != fiber.StatusOK {
		t.Errorf("invalid response")
	}

	rq2 := httptest.NewRequest("GET", "/?i=000000000000&m=pv_ins", nil)
	rq2.Header.Set("X-Real-IP", fake.IPv4Address())
	rq2.Header.Set(fiber.HeaderReferer, "http://example.net/page1")
	rs2, _ := app.Test(rq2)

	if rs2.StatusCode != fiber.StatusOK {
		t.Errorf("invalid response")
	}
}
func TestHTTP001(t *testing.T) {
	conf := newConfig("error", 0, "http://127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := getTestProjects()
	refererParser := newRefererParser()
	userAgentParser := newUserAgentParser()
	redisClient, redisClientErr := redisGetClient("redis://127.0.0.1:6379/0")
	if redisClientErr != nil {
		t.Error(redisClientErr)
		return
	}

	app := newHTTPServer(
		conf,
		geoParser,
		refererParser,
		userAgentParser,
		projectsManager,
		redisClient,
	)

	rq0 := httptest.NewRequest("GET", "/?i=000000000000&m=err_l&err=Sample&u=http%3A//example-not.net/", nil)
	rq0.Header.Set("X-Real-IP", fake.IPv4Address())
	rs0, _ := app.Test(rq0)

	if rs0.StatusCode != fiber.StatusForbidden {
		t.Errorf("invalid response")
	}

	rq1 := httptest.NewRequest("GET", "/?i=000000000000&m=err_l&err=Sample&u=http%3A//example.net/", nil)
	rq1.Header.Set("X-Real-IP", fake.IPv4Address())
	rs1, _ := app.Test(rq1)

	if rs1.StatusCode != fiber.StatusOK {
		t.Errorf("invalid response")
	}

	rq2 := httptest.NewRequest("GET", "/?i=000000000000&m=pv_ins", nil)
	rq2.Header.Set("X-Real-IP", fake.IPv4Address())
	rq2.Header.Set(fiber.HeaderReferer, "http://example.net/page1")
	rs2, _ := app.Test(rq2)

	if rs2.StatusCode != fiber.StatusOK {
		t.Errorf("invalid response")
	}
}
func TestHTTP2(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	conf := newConfig("error", 0, "http://127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	j, _ := projectsLoadJSON("./projects.json")
	projectsManager.load(j)
	redisClient, _ := redisGetClient("redis://127.0.0.1:6379/0")
	refererParser := newRefererParser()
	userAgentParser := newUserAgentParser()

	app := newHTTPServer(
		conf,
		geoParser,
		refererParser,
		userAgentParser,
		projectsManager,
		redisClient,
	)

	f, _ := fixtureLoad("./fixture.yml")
	r := f.record(refererParser, userAgentParser)

	recordCount := 0

	rq1 := httptest.NewRequest("GET", fmt.Sprintf(
		"/?m=pv_ins&i=%s",
		r.PublicInstanceID,
	), nil)

	rq1.Header.Set(fiber.HeaderXForwardedFor, r.IP.String())
	rq1.Header.Set(fiber.HeaderUserAgent, r.UserAgentResult.UaFull)
	rs1, _ := app.Test(rq1)

	if rs1.StatusCode != fiber.StatusFailedDependency {
		t.Errorf("invalid response")
	}

	for i := 0; i < 5000; i++ {

		r := f.record(refererParser, userAgentParser)
		if r.Mode > 99 {
			continue
		}

		recordCount++
		rq2 := httptest.NewRequest("GET", fmt.Sprintf(
			"/?m=pv_ins&i=%s",
			r.PublicInstanceID,
		), nil)

		rq2.Header.Set(fiber.HeaderXForwardedFor, r.IP.String())
		rq2.Header.Set(fiber.HeaderUserAgent, r.UserAgentResult.UaFull)
		rq2.Header.Set(fiber.HeaderReferer, r.PURL)
		rs2, _ := app.Test(rq2)

		if rs2.StatusCode != fiber.StatusOK {
			t.Errorf("invalid response")
		}

		recordCount++
		rq3 := httptest.NewRequest("GET", fmt.Sprintf(
			"/?m=pv_ins&i=%s&u=%s&cu=%s&t=%s&l=%s&ei=%s&em=%s&et=%s",
			r.PublicInstanceID,
			url.QueryEscape(r.PURL),
			url.QueryEscape(r.PCanonicalURL),
			url.QueryEscape(r.PTitle),
			url.QueryEscape(r.PLang),
			url.QueryEscape(r.PEntityID),
			url.QueryEscape(r.PEntityModule),
			url.QueryEscape(strconv.Itoa(int(r.PEntityTaxonomyID))),
		), nil)

		rq3.Header.Set(fiber.HeaderXForwardedFor, r.IP.String())
		rq3.Header.Set(fiber.HeaderUserAgent, r.UserAgentResult.UaFull)
		rq3.Header.Set(fiber.HeaderReferer, r.PURL)
		rs3, _ := app.Test(rq3)

		if rs3.StatusCode != fiber.StatusOK {
			t.Errorf("invalid response")
		}

	}

	time.Sleep(time.Duration(1) * time.Second)
}

func TestHTTP3(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	conf := newConfig("error", 0, "http://127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	j, _ := projectsLoadJSON("./projects.json")
	projectsManager.load(j)
	redisClient, _ := redisGetClient("redis://127.0.0.1:6379/0")
	refererParser := newRefererParser()
	userAgentParser := newUserAgentParser()

	app := newHTTPServer(
		conf,
		geoParser,
		refererParser,
		userAgentParser,
		projectsManager,
		redisClient,
	)

	f, _ := fixtureLoad("./fixture.yml")

	r := f.record(refererParser, userAgentParser)

	b, _ := json.Marshal(1)

	rq1 := httptest.NewRequest("POST", fmt.Sprintf("/?m=err&i=%s&u=%s", r.PublicInstanceID, r.PURL), strings.NewReader(string(b)))
	rq1.Header.Set(fiber.HeaderXForwardedFor, r.IP.String())
	rq1.Header.Set(fiber.HeaderUserAgent, r.UserAgentResult.UaFull)
	rq1.Header.Set(fiber.HeaderContentType, "application/json")
	rs1, _ := app.Test(rq1)

	if rs1.StatusCode != fiber.StatusUnprocessableEntity {
		t.Errorf("invalid response")
	}

	clientErrorCount := 0

	for i := 0; i < 50; i++ {

		r := f.record(refererParser, userAgentParser)
		if r.Mode > 99 {
			continue
		}

		postData := postRequest{
			ClientErrorMessage: "msg",
			ClientErrorObject:  "errObject",
		}

		b, _ := json.Marshal(postData)

		rq1 := httptest.NewRequest("POST", fmt.Sprintf("/?m=err&i=%s&u=%s", r.PublicInstanceID, r.PURL), strings.NewReader(string(b)))
		rq1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
		rq1.Header.Set(fiber.HeaderUserAgent, r.UserAgentResult.UaFull)
		rq1.Header.Set(fiber.HeaderContentType, "application/json")
		rs1, _ := app.Test(rq1)

		if rs1.StatusCode != fiber.StatusOK {
			t.Errorf("invalid response")
		}

		clientErrorCount++
	}

	time.Sleep(time.Duration(1) * time.Second)
}

func TestHTTP4(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	conf := newConfig("error", 0, "http://127.0.0.1")
	geoParser := getGeoParser()
	projectsManager := newProjectsManager()
	j, _ := projectsLoadJSON("./projects.json")
	projectsManager.load(j)
	redisClient, _ := redisGetClient("redis://127.0.0.1:6379/0")
	refererParser := newRefererParser()
	userAgentParser := newUserAgentParser()

	app := newHTTPServer(
		conf,
		geoParser,
		refererParser,
		userAgentParser,
		projectsManager,
		redisClient,
	)

	f, _ := fixtureLoad("./fixture.yml")

	for i := 0; i < 5000; i++ {

		r := f.record(refererParser, userAgentParser)
		if r.Mode < 100 || r.EventCount == 0 {
			continue
		}

		ev := make([]postRequestEvent, 0)

		for _, e := range r.Events {
			e := postRequestEvent{
				Category: e.ECategory,
				Action:   e.EAction,
				Label:    e.ELabel,
				Ident:    e.EIdent,
				Value:    int64(e.EValue),
			}

			ev = append(ev, e)
		}

		apiData := postRequestAPI{
			ClientIP:           r.IP.String(),
			ClientUserAgent:    r.UserAgentResult.UaFull,
			PrivateInstanceKey: "000000000000111111111111",
		}

		if f.rand(1, 5) == 3 {
			ut := time.Now().Unix() + int64(f.rand(-86400, 86400))
			apiData.ClientTime = ut
		}

		postData := postRequest{
			API:    &apiData,
			Events: &ev,
		}

		b, _ := json.Marshal(postData)

		rq1 := httptest.NewRequest("POST", fmt.Sprintf("/?m=e_api&i=%s", r.PublicInstanceID), strings.NewReader(string(b)))
		rq1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
		rq1.Header.Set(fiber.HeaderUserAgent, r.UserAgentResult.UaFull)
		rq1.Header.Set(fiber.HeaderContentType, "application/json")
		rs1, _ := app.Test(rq1)

		if rs1.StatusCode != fiber.StatusOK {
			t.Errorf("invalid response")
		}

		apiDataInvalid := postRequestAPI{}

		if fake.Number(1, 3) == 1 {
			apiDataInvalid.ClientIP = fake.IPv4Address()
		} else if fake.Number(1, 3) == 1 {
			apiDataInvalid.ClientUserAgent = r.UserAgentResult.UaFull
		} else if fake.Number(1, 3) == 1 {
			apiDataInvalid.PrivateInstanceKey = "000000000000111111111111"
		}

		postDataInvalid := postRequest{
			API:    &apiDataInvalid,
			Events: &ev,
		}

		bInvalid, _ := json.Marshal(postDataInvalid)

		rq2 := httptest.NewRequest("POST", fmt.Sprintf("/?m=e_api&i=%s", r.PublicInstanceID), strings.NewReader(string(bInvalid)))
		rq2.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
		rq2.Header.Set(fiber.HeaderUserAgent, r.UserAgentResult.UaFull)
		rq2.Header.Set(fiber.HeaderContentType, "application/json")
		rs2, _ := app.Test(rq2)

		if rs2.StatusCode == fiber.StatusOK {
			t.Errorf("invalid response")
		}
	}
}
