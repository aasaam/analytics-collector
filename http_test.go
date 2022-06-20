package main

const sampleUserAgent = "Mozilla/5.0 (Linux; Android 5.0.2; SAMSUNG SM-A500FU Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/3.3 Chrome/38.0.2125.102 Mobile Safari/537.36"

// func TestHTTP1(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "")
// 	geoParser := getGeoParser()
// 	projectsManager := newProjectsManager()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	statics := []string{
// 		"/a.js",
// 		"/l.js",
// 	}

// 	for _, p := range statics {
// 		r := httptest.NewRequest("GET", p, nil)
// 		rs, _ := app.Test(r)

// 		if rs.StatusCode != fiber.StatusOK {
// 			t.Errorf("invalid response")
// 		}
// 	}
// }
// func TestHTTP11(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "")
// 	geoParser := getGeoParser()
// 	projectsManager := newProjectsManager()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	statics := []string{
// 		"/a.js",
// 		"/a.src.js",
// 		"/l.src.js",
// 		"/amp.json",
// 	}

// 	for _, p := range statics {
// 		r := httptest.NewRequest("GET", p, nil)
// 		rs, _ := app.Test(r)

// 		if rs.StatusCode != fiber.StatusOK {
// 			t.Errorf("invalid response")
// 		}
// 	}
// }

// func TestHTTP110(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "")
// 	geoParser := getGeoParser()
// 	projectsManager := newProjectsManager()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	statics := []string{
// 		"/a.js?q",
// 		"/a.src.js?q",
// 		"/l.src.js?foo=1",
// 		"/amp.json?q",
// 	}

// 	for _, p := range statics {
// 		r := httptest.NewRequest("GET", p, nil)
// 		rs, _ := app.Test(r)
// 		if rs.StatusCode != fiber.StatusForbidden {
// 			t.Errorf("invalid response")
// 		}
// 	}
// }

// func TestHTTP2(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "")
// 	geoParser := getGeoParser()
// 	projectsManager := newProjectsManager()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	r404 := []string{
// 		"/favicon.ico",
// 		"/ensure-not-exist",
// 	}

// 	for _, p := range r404 {
// 		r := httptest.NewRequest("GET", p, nil)
// 		rs, _ := app.Test(r)

// 		if rs.StatusCode < 400 && rs.StatusCode >= 500 {
// 			t.Errorf("invalid response")
// 		}
// 	}

// }
// func TestHTTP3(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
// 	geoParser := getGeoParser()
// 	projectsManager := newProjectsManager()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	r1 := httptest.NewRequest("GET", "/metrics", nil)
// 	r1.Header.Set(fiber.HeaderXForwardedFor, "192.168.1.1")
// 	rs1, _ := app.Test(r1)

// 	if rs1.StatusCode == fiber.StatusOK {
// 		t.Errorf("invalid response")
// 	}

// 	r2 := httptest.NewRequest("GET", "/metrics", nil)
// 	r2.Header.Set(fiber.HeaderXForwardedFor, "127.0.0.1")
// 	rs2, _ := app.Test(r2)

// 	if rs2.StatusCode != fiber.StatusOK {
// 		t.Errorf("invalid response")
// 	}
// }
// func TestHTTP10(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
// 	geoParser := getGeoParser()
// 	projectsManager := newProjectsManager()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	r1 := httptest.NewRequest("PATCH", "/?m=pv_ins&i=000000000000", nil)
// 	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
// 	rs1, _ := app.Test(r1)

// 	if rs1.StatusCode == fiber.StatusOK {
// 		t.Errorf("invalid response")
// 	}
// }
// func TestHTTP12(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
// 	geoParser := getGeoParser()
// 	projectsManager := newProjectsManager()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	r1 := httptest.NewRequest("PATCH", "/?m=pv_ins&i=00000000000_", nil)
// 	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
// 	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
// 	rs1, _ := app.Test(r1)

// 	if rs1.StatusCode == fiber.StatusOK {
// 		t.Errorf("invalid response")
// 	}
// }
// func TestHTTP20(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
// 	geoParser := getGeoParser()
// 	projectsManager := getTestProjects()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	r1 := httptest.NewRequest("GET", "/?m=pv_ins&i=000000000000&u=https%3A%2F%2Fexample.com%2F", nil)
// 	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
// 	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
// 	rs1, _ := app.Test(r1)

// 	if rs1.StatusCode != fiber.StatusOK {
// 		t.Errorf("invalid response")
// 	}
// }
// func TestHTTP21(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
// 	geoParser := getGeoParser()
// 	projectsManager := getTestProjects()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	r1 := httptest.NewRequest("GET", "/?m=pv_ins&i=000000000000", nil)
// 	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
// 	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
// 	r1.Header.Set(fiber.HeaderReferer, "http://example.com")
// 	rs1, _ := app.Test(r1)

// 	if rs1.StatusCode != fiber.StatusOK {
// 		t.Errorf("invalid response")
// 	}
// }
// func TestHTTP22(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
// 	geoParser := getGeoParser()
// 	projectsManager := getTestProjects()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	r1 := httptest.NewRequest("GET", "/?m=pv_ins&i=000000000000", nil)
// 	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
// 	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
// 	rs1, _ := app.Test(r1)

// 	if rs1.StatusCode == fiber.StatusOK {
// 		t.Errorf("invalid response")
// 	}
// }

// func TestHTTP30(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
// 	geoParser := getGeoParser()
// 	projectsManager := getTestProjects()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	r1 := httptest.NewRequest("POST", "/?m=err&i=000000000000&u=https%3A%2F%2Fexample.com%2F", strings.NewReader(`{"foo":true"}`))
// 	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
// 	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
// 	r1.Header.Set(fiber.HeaderContentType, "application/json")
// 	rs1, _ := app.Test(r1)

// 	if rs1.StatusCode == fiber.StatusOK {
// 		t.Errorf("invalid response")
// 	}
// }

// func TestHTTP31(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
// 	geoParser := getGeoParser()
// 	projectsManager := getTestProjects()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	postData := postRequest{
// 		ClientErrorMessage: "msg",
// 		ClientErrorObject:  "errObject",
// 	}

// 	b, _ := json.Marshal(postData)

// 	r1 := httptest.NewRequest("POST", "/?m=err&i=000000000000&u=https%3A%2F%2Fexample.com%2F", strings.NewReader(string(b)))
// 	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
// 	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
// 	r1.Header.Set(fiber.HeaderContentType, "application/json")
// 	rs1, _ := app.Test(r1)

// 	if rs1.StatusCode != fiber.StatusOK {
// 		t.Errorf("invalid response")
// 	}
// }
// func TestHTTP32(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
// 	geoParser := getGeoParser()
// 	projectsManager := getTestProjects()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	postData := postRequest{}

// 	b, _ := json.Marshal(postData)

// 	r1 := httptest.NewRequest("POST", "/?m=e_api&i=000000000000", strings.NewReader(string(b)))
// 	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
// 	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
// 	r1.Header.Set(fiber.HeaderContentType, "application/json")
// 	rs1, _ := app.Test(r1)

// 	if rs1.StatusCode == fiber.StatusOK {
// 		t.Errorf("invalid response")
// 	}
// }
// func TestHTTP33(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
// 	geoParser := getGeoParser()
// 	projectsManager := getTestProjects()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	api := postRequestAPI{
// 		PrivateInstanceKey: "000000000000111111111111",
// 		ClientIP:           "8.8.8.8",
// 		ClientUserAgent:    "curl 1.1.2",
// 		ClientTime:         time.Now().Unix(),
// 	}

// 	ev1 := postRequestEvent{
// 		Category: "cat",
// 		Action:   "act",
// 		Label:    "lab",
// 		Value:    1,
// 	}
// 	ev2 := postRequestEvent{
// 		Category: "!@#",
// 		Action:   "!@#",
// 		Label:    "lab",
// 		Value:    1,
// 	}

// 	events := []postRequestEvent{ev1, ev2}

// 	postData := postRequest{
// 		API:    &api,
// 		Events: &events,
// 	}

// 	b, _ := json.Marshal(postData)

// 	r1 := httptest.NewRequest("POST", "/?m=e_api&i=000000000000", strings.NewReader(string(b)))
// 	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
// 	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
// 	r1.Header.Set(fiber.HeaderContentType, "application/json")
// 	rs1, _ := app.Test(r1)

// 	if rs1.StatusCode != fiber.StatusOK {
// 		t.Errorf("invalid response")
// 	}
// }
// func TestHTTP40(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 		return
// 	}

// 	c1 := newConfig("error", 0, true, "http://127.0.0.1", "127.0.0.1")
// 	geoParser := getGeoParser()
// 	projectsManager := getTestProjects()

// 	app := newHTTPServer(c1, geoParser, projectsManager)

// 	jsonSample := `{"cid_std":"MTY0OTE4MzAzMzoxNjQ5MTgzMDMzOjFmdnRmZzF2OWdmMTE1YnI=","p":{"u":"https://www.example.net/%D8%A8%D8%AE%D8%B4-%D8%B1%D8%B3%D8%A7%D9%86%D9%87-71/1011644-%D9%BE%DB%8C%D8%A7%D9%85-%D8%B1%D9%87%D8%A8%D8%B1-%D9%85%D8%B9%D8%B8%D9%85-%D8%A7%D9%86%D9%82%D9%84%D8%A7%D8%A8-%D8%A8%D9%87-%D9%85%D9%86%D8%A7%D8%B3%D8%A8%D8%AA-%D8%A2%D8%BA%D8%A7%D8%B2-%D8%B3%D8%A7%D9%84-%D9%86%D9%88%DB%8C%D8%AF-%D8%A7%D9%85%DB%8C%D8%AF%D8%A8%D8%AE%D8%B4-%D8%A7%D9%82%D8%AA%D8%B5%D8%A7%D8%AF%DB%8C-%D8%A7%D8%B2-%D8%B7%D8%B1%D9%81-%D8%B1%D9%87%D8%A8%D8%B1-%D8%A7%D9%86%D9%82%D9%84%D8%A7%D8%A8-%D8%A8%D8%B1%D8%A7%DB%8C-%D8%B3%D8%A7%D9%84-%D9%82%D8%B1%D9%86-%D8%AC%D8%AF%DB%8C%D8%AF","t":"پیام رهبر معظم انقلاب به مناسبت آغاز سال 1401 | نوید امیدبخش اقتصادی از طرف رهبر انقلاب برای سال و قرن جدید","l":"fa","cu":"https://www.example.net/بخش-%D8%B1%D8%B3%D8%A7%D9%86%D9%87-71/1011644-%D9%BE%DB%8C%D8%A7%D9%85-%D8%B1%D9%87%D8%A8%D8%B1-%D9%85%D8%B9%D8%B8%D9%85-%D8%A7%D9%86%D9%82%D9%84%D8%A7%D8%A8-%D8%A8%D9%87-%D9%85%D9%86%D8%A7%D8%B3%D8%A8%D8%AA-%D8%A2%D8%BA%D8%A7%D8%B2-%D8%B3%D8%A7%D9%84-%D9%86%D9%88%DB%8C%D8%AF-%D8%A7%D9%85%DB%8C%D8%AF%D8%A8%D8%AE%D8%B4-%D8%A7%D9%82%D8%AA%D8%B5%D8%A7%D8%AF%DB%8C-%D8%A7%D8%B2-%D8%B7%D8%B1%D9%81-%D8%B1%D9%87%D8%A8%D8%B1-%D8%A7%D9%86%D9%82%D9%84%D8%A7%D8%A8-%D8%A8%D8%B1%D8%A7%DB%8C-%D8%B3%D8%A7%D9%84-%D9%82%D8%B1%D9%86-%D8%AC%D8%AF%DB%8C%D8%AF","ei":"","em":"","r":"","bc":{},"scr":"1920x1080","vps":"1868x344","cd":"24","k":"حضرت آیت الله العظمی خامنه ای,پیام نوروزی رهبر انقلاب,پیام نوروزی رهبر معظم انقلاب,پیام آغاز سال 1401 رهبر انقلاب","rs":"","dpr":"1","if":false,"ts":false,"sot":"landscape-primary","prf":{"dlt":"0","tct":"66","srt":"179","pdt":"93","rt":"0","dit":"617","clt":"617","r":38}}}`

// 	r1 := httptest.NewRequest("POST", "/?m=pv_js&i=000000000000", strings.NewReader(jsonSample))
// 	r1.Header.Set(fiber.HeaderXForwardedFor, "1.1.1.1")
// 	r1.Header.Set(fiber.HeaderUserAgent, sampleUserAgent)
// 	r1.Header.Set(fiber.HeaderContentType, "text/plain;charset=UTF-8")
// 	rs1, _ := app.Test(r1)

// 	if rs1.StatusCode != fiber.StatusOK {
// 		t.Errorf("invalid response")
// 	}
// }
