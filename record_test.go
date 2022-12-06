package main

import (
	"encoding/base64"
	"net"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestRecord1(t *testing.T) {
	_, e := validateMode("ensure-not-exist")
	if e == nil {
		t.Errorf("invalid mode of record")
	}
}
func TestRecord01(t *testing.T) {
	r := record{}
	_, finalizeErr := r.finalize()
	if finalizeErr == nil {
		t.Error("must invalid record")
	}

}
func TestRecord2(t *testing.T) {
	_, e1 := newRecord("invalid", "000000000000")
	if e1 == nil {
		t.Errorf("invalid init record")
	}
	_, e2 := newRecord("pv_js", "00000000000!")
	if e2 == nil {
		t.Errorf("invalid init record")
	}
	recordSample, e3 := newRecord("pv_js", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	recordSample.CID = clientIDNoneSTD([]string{"amp"}, clientIDTypeOther)

	projects := projects{}

	em1 := recordSample.verify(&projects, "")
	if em1 != &errorRecordNotProcessedYet {
		t.Errorf("invalid verify")
	}

	_, finalizeErr := recordSample.finalize()
	if finalizeErr != nil {
		t.Error(finalizeErr)
	}

	if !recordSample.isPageView() || recordSample.isImage() {
		t.Errorf("invalid is page view")
	}

	ev1, e3 := newRecord("e_js_pv", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	if ev1.isPageView() || ev1.isImage() {
		t.Errorf("invalid is page view")
	}
}
func TestRecord3(t *testing.T) {
	refererParser := newRefererParser()
	recordSample, e3 := newRecord("pv_il", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	recordSample.CID = clientIDNoneSTD([]string{"amp"}, clientIDTypeOther)

	recordSample.setQueryParameters(
		"https://www.example.com/?utm_source=source&utm_medium=medium&utm_campaign=sale1&utm_id=id&utm_term=keyword1&utm_content=content",
		"http://example.com",
		"Page title",
		"en",
		"0",
		"home",
		"R0000",
	)

	recordSample.setReferer(refererParser, getURL("http://www.google.com"))

	if recordSample.PRefererURL.RefName != "Google" {
		t.Errorf("invalid referer")
	}

	_, finalizeErr := recordSample.finalize()
	if finalizeErr != nil {
		t.Error(finalizeErr)
	}

	if recordSample.Utm.UtmSource != "source" {
		t.Errorf("invalid utm")
	}

	if !recordSample.isPageView() || !recordSample.isImage() {
		t.Errorf("invalid is page view")
	}

}
func TestRecord4(t *testing.T) {
	recordSample, e3 := newRecord("pv_js", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	pm := getTestProjects()

	if recordSample.verify(pm, "") != &errorRecordNotProcessedYet {
		t.Errorf("invalid verify")
	}

	recordSample.CID = clientIDNoneSTD([]string{"amp"}, clientIDTypeOther)

	recordSample.IP = net.ParseIP("1.1.1.1")

	if recordSample.verify(pm, "") != &errorURLRequiredAndMustBeValid {
		t.Errorf("invalid verify")
	}

	u1, _ := url.Parse("http://not-example.com")

	recordSample.PURL = u1.String()
	recordSample.pURL = u1

	if recordSample.verify(pm, "") != &errorProjectPublicIDAndURLDidNotMatched {
		t.Errorf("invalid verify")
	}

	u2, _ := url.Parse("http://example.com")

	recordSample.PURL = u2.String()
	recordSample.pURL = u2

	if recordSample.verify(pm, "") != nil {
		t.Errorf("invalid verify")
	}
}
func TestRecord5(t *testing.T) {
	recordSample, e3 := newRecord("e_api", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	pm := getTestProjects()
	recordSample.IP = net.ParseIP("1.1.1.1")
	u2, _ := url.Parse("http://example.com")
	recordSample.PURL = u2.String()

	if recordSample.verify(pm, "") != &errorAPIPrivateKeyFailed {
		t.Errorf("invalid verify")
	}
}

func TestRecord10(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	geoParser := getGeoParser()
	refererParser := newRefererParser()

	recordSample, e3 := newRecord("e_js_pv", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	pm := getTestProjects()

	recordSample.IP = net.ParseIP("1.1.1.1")

	u2, _ := url.Parse("http://example.com")

	recordSample.PURL = u2.String()

	ev1 := postRequestEvent{
		Category: "cat",
		Action:   "act",
		Label:    "lab",
		Ident:    "100",
		Value:    1,
	}
	ev2 := postRequestEvent{
		Category: "!@#",
		Action:   "!@#",
		Label:    "lab",
		Ident:    "100",
		Value:    1,
	}

	postPage := postRequestPage{
		URL: "http://example.com",
	}

	events := []postRequestEvent{ev1, ev2}

	initTime := time.Now().Add(time.Duration(-60) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	validCID := base64.StdEncoding.EncodeToString([]byte(cid))

	post := postRequest{
		CIDStd: validCID,
		Page:   &postPage,
		Events: &events,
	}

	recordSample.setPostRequest(&post, refererParser, geoParser)

	if recordSample.verify(pm, "") != nil {
		t.Errorf("invalid verify")
	}
}
func TestRecord11(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	geoParser := getGeoParser()
	refererParser := newRefererParser()

	recordSample, e3 := newRecord("e_js_pv", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	pm := getTestProjects()

	recordSample.IP = net.ParseIP("1.1.1.1")

	u2, _ := url.Parse("http://example.com")

	recordSample.PURL = u2.String()

	ev1 := postRequestEvent{
		Category: "cat",
		Action:   "act",
		Label:    "lab",
		Ident:    "100",
		Value:    1,
	}
	ev2 := postRequestEvent{
		Category: "!@#",
		Action:   "!@#",
		Label:    "lab",
		Ident:    "100",
		Value:    1,
	}

	bc := postRequestBreadcrumb{
		N1: "name 1",
		N2: "name 2",
		N3: "name 3",
		N4: "name 4",
		N5: "name 5",
		U1: "http://example.com/path-1",
		U2: "http://example.com/path-1/path-2",
		U3: "http://example.com/path-1/path-2/path-3",
		U4: "http://example.com/path-1/path-2/path-3/path-4",
		U5: "http://example.com/path-1/path-2/path-3/path-4/path-5",
	}

	sg := postRequestSegment{
		S1N: "home",
		S2N: "product",
		S3N: "top",
		S4N: "special",
		S5N: "camera1",
		S1V: "fa",
		S2V: "1",
		S3V: "2",
		S4V: "3",
		S5V: "4",
	}

	postPage := postRequestPage{
		URL:                  "http://example.com",
		PageBreadcrumbObject: &bc,
		Seg:                  &sg,
	}

	events := []postRequestEvent{ev1, ev2}

	cid := "aa:aa:0000000000000000"
	inValidCID := base64.StdEncoding.EncodeToString([]byte(cid))

	post := postRequest{
		CIDStd: inValidCID,
		Page:   &postPage,
		Events: &events,
	}

	recordSample.setPostRequest(&post, refererParser, geoParser)

	if recordSample.verify(pm, "") != nil {
		t.Errorf("invalid verify")
	}
}

func TestRecord12(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	geoParser := getGeoParser()
	refererParser := newRefererParser()

	recordSample, e3 := newRecord("e_js_pv", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	recordSample.CID = clientIDNoneSTD([]string{"amp"}, clientIDTypeOther)

	pm := getTestProjects()

	recordSample.IP = net.ParseIP("1.1.1.1")

	u2, _ := url.Parse("http://example-another.com")

	recordSample.PURL = u2.String()

	ev1 := postRequestEvent{
		Category: "cat",
		Action:   "act",
		Label:    "lab",
		Ident:    "100",
		Value:    1,
	}
	ev2 := postRequestEvent{
		Category: "!@#",
		Action:   "!@#",
		Label:    "lab",
		Ident:    "100",
		Value:    1,
	}

	postPage := postRequestPage{}

	events := []postRequestEvent{ev1, ev2}

	initTime := time.Now().Add(time.Duration(-60) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	validCID := base64.StdEncoding.EncodeToString([]byte(cid))

	post := postRequest{
		CIDStd: validCID,
		Page:   &postPage,
		Events: &events,
	}

	recordSample.setPostRequest(&post, refererParser, geoParser)

	if recordSample.verify(pm, "") != &errorProjectPublicIDAndURLDidNotMatched {
		t.Errorf("invalid verify")
	}
}
func TestRecord131(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	geoParser := getGeoParser()
	userAgentParser := newUserAgentParser()
	pm := getTestProjects()

	recordSample, e3 := newRecord("e_api", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	recordSample.IP = net.ParseIP("1.1.1.1")

	u2, _ := url.Parse("http://example.net")

	recordSample.PURL = u2.String()

	postPage := postRequestPage{
		URL: "http://example.net/",
	}

	events := []postRequestEvent{}

	initTime := time.Now().Add(time.Duration(-60) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	validCID := base64.StdEncoding.EncodeToString([]byte(cid))

	post := postRequest{
		CIDStd: validCID,
		Page:   &postPage,
		Events: &events,
		API: &postRequestAPI{
			PrivateInstanceKey: "000000000000111111111111",
			ClientIP:           "127.1.1.1",
			ClientTime:         time.Now().Unix(),
		},
	}

	e := recordSample.setAPI(pm, userAgentParser, geoParser, &post)

	if e != &errorAPIClientUserAgentNotValid {
		t.Errorf("invalid verify")
	}
}
func TestRecord123(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	geoParser := getGeoParser()
	userAgentParser := newUserAgentParser()
	pm := getTestProjects()

	recordSample, e3 := newRecord("e_api", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	recordSample.IP = net.ParseIP("1.1.1.1")

	u2, _ := url.Parse("http://example.net")

	recordSample.PURL = u2.String()

	postPage := postRequestPage{
		URL: "http://example.net/",
	}

	events := []postRequestEvent{}

	initTime := time.Now().Add(time.Duration(-60) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	validCID := base64.StdEncoding.EncodeToString([]byte(cid))

	post := postRequest{
		CIDStd: validCID,
		Page:   &postPage,
		Events: &events,
		API: &postRequestAPI{
			PrivateInstanceKey: "invalid",
			ClientIP:           "127.1.1.1",
			ClientUserAgent:    "AdsBot-Google (+http://www.google.com/adsbot.html)",
			ClientTime:         time.Now().Unix(),
		},
	}

	e := recordSample.setAPI(pm, userAgentParser, geoParser, &post)

	if e != &errorAPIPrivateKeyFailed {
		t.Errorf("invalid verify")
	}
}
func TestRecord125(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	geoParser := getGeoParser()
	userAgentParser := newUserAgentParser()
	pm := getTestProjects()

	recordSample, e3 := newRecord("e_api", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	recordSample.IP = net.ParseIP("1.1.1.1")

	u2, _ := url.Parse("http://example.net")

	recordSample.PURL = u2.String()

	postPage := postRequestPage{
		URL: "http://example.net/",
	}

	events := []postRequestEvent{}

	initTime := time.Now().Add(time.Duration(-60) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	validCID := base64.StdEncoding.EncodeToString([]byte(cid))

	post := postRequest{
		CIDStd: validCID,
		Page:   &postPage,
		Events: &events,
		API: &postRequestAPI{
			PrivateInstanceKey: "000000000000111111111111",
			ClientUserAgent:    "AdsBot-Google (+http://www.google.com/adsbot.html)",
			ClientTime:         time.Now().Unix(),
		},
	}

	e := recordSample.setAPI(pm, userAgentParser, geoParser, &post)

	if e != &errorAPIClientIPNotValid {
		t.Errorf("invalid verify")
	}
}
func TestRecord140(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	geoParser := getGeoParser()
	userAgentParser := newUserAgentParser()
	pm := getTestProjects()

	recordSample, e3 := newRecord("e_api", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	recordSample.IP = net.ParseIP("1.1.1.1")

	u2, _ := url.Parse("http://example.net")

	recordSample.PURL = u2.String()

	postPage := postRequestPage{
		URL: "http://example.net/",
	}

	events := []postRequestEvent{}

	initTime := time.Now().Add(time.Duration(-60) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	validCID := base64.StdEncoding.EncodeToString([]byte(cid))

	post := postRequest{
		CIDStd: validCID,
		Page:   &postPage,
		Events: &events,
		API: &postRequestAPI{
			PrivateInstanceKey: "000000000000111111111111",
			ClientIP:           "1.1.1.1",
			ClientUserAgent:    "AdsBot-Google (+http://www.google.com/adsbot.html)",
			ClientTime:         time.Now().Unix(),
		},
	}

	e := recordSample.setAPI(pm, userAgentParser, geoParser, &post)

	if e != nil {
		t.Error(e)
	}
}
func TestRecord124(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	geoParser := getGeoParser()
	userAgentParser := newUserAgentParser()
	pm := getTestProjects()

	recordSample, e3 := newRecord("e_api", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	recordSample.IP = net.ParseIP("1.1.1.1")

	u2, _ := url.Parse("http://example.net")

	recordSample.PURL = u2.String()

	postPage := postRequestPage{
		URL: "http://example.net/",
	}

	events := []postRequestEvent{}

	initTime := time.Now().Add(time.Duration(-60) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	validCID := base64.StdEncoding.EncodeToString([]byte(cid))

	post := postRequest{
		CIDStd: validCID,
		Page:   &postPage,
		Events: &events,
	}

	e := recordSample.setAPI(pm, userAgentParser, geoParser, &post)

	if e != &errorAPIFieldsAreMissing {
		t.Errorf("invalid verify")
	}
}
func TestRecord121(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	geoParser := getGeoParser()
	refererParser := newRefererParser()

	recordSample, e3 := newRecord("e_js_pv", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	recordSample.CID = clientIDNoneSTD([]string{"amp"}, clientIDTypeOther)

	pm := getTestProjects()

	recordSample.IP = net.ParseIP("1.1.1.1")

	u2, _ := url.Parse("http://example.net")

	recordSample.PURL = u2.String()

	postPage := postRequestPage{
		URL: "http://example.net/",
	}

	events := []postRequestEvent{}

	initTime := time.Now().Add(time.Duration(-60) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	validCID := base64.StdEncoding.EncodeToString([]byte(cid))

	post := postRequest{
		CIDStd: validCID,
		Page:   &postPage,
		Events: &events,
	}

	recordSample.setPostRequest(&post, refererParser, geoParser)

	if recordSample.verify(pm, "") != &errorEventsAreEmpty {
		t.Errorf("invalid verify")
	}
}

func TestRecord13(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	geoParser := getGeoParser()
	refererParser := newRefererParser()

	recordSample, e3 := newRecord("err", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	recordSample.CID = clientIDNoneSTD([]string{"amp"}, clientIDTypeOther)

	pm := getTestProjects()

	recordSample.IP = net.ParseIP("1.1.1.1")

	u2, _ := url.Parse("http://not-example.com")

	recordSample.PURL = u2.String()

	perf := postRequestPerformanceData{
		PerfPageLoadTime: "100",
	}

	bc := postRequestBreadcrumb{
		N1: "name 1",
		U1: "http://not-example.com/path-1",
	}

	geo := postRequestGeographyData{
		Lat: 90,
		Lon: -180,
	}

	postPage := postRequestPage{
		URL:                   "http://not-example.com",
		RefererURL:            "https://duckduckgo.com/",
		ScreenSize:            "1024x768",
		ViewportSize:          "800x600",
		DevicePixelRatio:      "1.1",
		ColorDepth:            "24",
		ScreenOrientationType: "l-p",
		RefererSessionURL:     "https://www.google.com",
		PerformanceData:       &perf,
		PageBreadcrumbObject:  &bc,
		GeographyData:         &geo,
	}

	initTime := time.Now().Add(time.Duration(-60) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	validCID := base64.StdEncoding.EncodeToString([]byte(cid))

	post := postRequest{
		ClientErrorMessage: "message",
		ClientErrorObject:  `{"foo":true}`,
		CIDStd:             validCID,
		Page:               &postPage,
	}

	recordSample.setPostRequest(&post, refererParser, geoParser)

	if recordSample.verify(pm, "") != &errorProjectPublicIDAndURLDidNotMatched {
		t.Errorf("invalid verify")
	}
}
func TestRecord14(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	geoParser := getGeoParser()
	refererParser := newRefererParser()

	recordSample, e3 := newRecord("pv_amp", "000000000000")
	if e3 != nil {
		t.Error(e3)
	}

	recordSample.CID = clientIDNoneSTD([]string{"amp"}, clientIDTypeOther)

	pm := getTestProjects()

	recordSample.IP = net.ParseIP("1.1.1.1")

	u2, _ := url.Parse("http://example.com")

	recordSample.PURL = u2.String()

	perf := postRequestPerformanceData{
		PerfPageLoadTime: "a200",
	}

	postPage := postRequestPage{
		URL:               "http://not-example.com",
		RefererURL:        "https://duckduckgo.com/",
		ScreenSize:        "1024x768",
		ViewportSize:      "800x600",
		RefererSessionURL: "https://www.google.com",
		PerformanceData:   &perf,
	}

	initTime := time.Now().Add(time.Duration(-60) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	validCID := base64.StdEncoding.EncodeToString([]byte(cid))

	post := postRequest{
		CIDAmp:             "ampCIDString",
		ClientErrorMessage: "message",
		ClientErrorObject:  `{"foo":true}`,
		CIDStd:             validCID,
		Page:               &postPage,
	}

	recordSample.setPostRequest(&post, refererParser, geoParser)

	if recordSample.verify(pm, "") != &errorProjectPublicIDAndURLDidNotMatched {
		t.Errorf("invalid verify")
	}
}
