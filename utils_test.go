package main

import (
	"testing"
	"time"
)

func TestSanitizeName(t *testing.T) {
	if sanitizeName("VIDEO") != "video" {
		t.Errorf("invalid sanitize")
	}

	if sanitizeName("!!!") != "" {
		t.Errorf("invalid sanitize")
	}
}

func TestBoolUint8(t *testing.T) {
	if boolUint8(true) != 1 {
		t.Errorf("invalid sanitize")
	}
	if boolUint8(false) != 0 {
		t.Errorf("invalid sanitize")
	}
}
func TestSanitizeEntityID(t *testing.T) {
	if sanitizeEntityID("1") != "1" {
		t.Errorf("invalid sanitize")
	}

	if sanitizeEntityID("!!!") != "" {
		t.Errorf("invalid sanitize")
	}
}
func TestSanitizeLanguage(t *testing.T) {
	if sanitizeLanguage("fa-IR") != "fa" {
		t.Errorf("invalid sanitize")
	}
	if sanitizeLanguage("per") != "fa" {
		t.Errorf("invalid sanitize")
	}

	if sanitizeLanguage("!!!") != "" {
		t.Errorf("invalid sanitize")
	}
}
func TestSanitizeEntityTaxonomyID(t *testing.T) {
	if sanitizeEntityTaxonomyID("G0000") != "G0000" {
		t.Errorf("invalid sanitize")
	}

	if sanitizeEntityTaxonomyID("!!!") != "" {
		t.Errorf("invalid sanitize")
	}
}

func TestHash(t *testing.T) {
	s := checksum("1")
	if len(s) != 24 {
		t.Errorf("invalid hash")
	}
}
func TestHash2(t *testing.T) {
	s := checksum("")
	if len(s) != 24 {
		t.Errorf("invalid hash")
	}
}
func TestIsValidURL(t *testing.T) {
	if sanitizeURL("\x18") != "" {
		t.Errorf("invalid url")
	}

	if isValidURL("\x18") {
		t.Errorf("invalid url")
	}

	if isValidURL("") {
		t.Errorf("invalid url")
	}

	if !isValidURL("http://google.com") {
		t.Errorf("invalid url")
	}
}

func TestParseKeywords(t *testing.T) {
	k1 := parseKeywords("1,2,3,4,5,6,7,8,9,10,11")
	if len(k1) != 10 {
		t.Errorf("invalid keyword parse")
	}

	k2 := parseKeywords("")

	if len(k2) != 0 {
		t.Errorf("invalid keyword parse")
	}

	k3 := parseKeywords("1,2")

	if len(k3) != 2 {
		t.Errorf("invalid keyword parse")
	}
}
func TestURLDomainParse(t *testing.T) {
	urls := []string{
		"https://www.xn--mgbtj4c7ad63e.xn--mgba3a4f16a",
		"https://subdomain.فروشگاه.ایران/",
		"https://sub.of.google.com/",
		"http://192.168.1.1/",
		"http://localhost/",
	}

	for _, ur := range urls {
		uu := getURL(ur)
		if uu == nil {
			t.Errorf("invaid url")
		}

		d1 := getDomain(uu)

		if d1 == "" {
			t.Errorf("invaid domain")
		}

		d2 := getSecondDomainLevel(uu)
		if d2 == "" {
			t.Errorf("invaid second level domain")
		}
	}

	if getURL("") != nil {
		t.Errorf("empty url string")
	}
}

func TestGetCursorID(t *testing.T) {
	c1, _ := getCursorID()
	time.Sleep(time.Duration(2) * time.Millisecond)
	c2, _ := getCursorID()
	if c1 == c2 {
		t.Errorf("must not same cursor")
	}
}

func TestGetURLPath(t *testing.T) {
	urls := []string{
		"https://www.google.com/path/تست/file.html?",
		"https://www.google.com/path/%D8%AA%D8%B3%D8%AA/file.html?",
		"https://www.google.com/path/to/file.html?foo=bar",
		"https://www.google.com/path/to/file.html?foo=bar#fragment",
		"https://www.google.com/path/to/file.html?foo=bar#/extra",
	}

	for _, u := range urls {
		p := getURLPath(getURL(u))
		s := getURLString(getURL(u))
		if p == "" || s == "" {
			t.Errorf("invalid url path")
		}
	}

}

func BenchmarkChecksum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		checksum("a")
	}
}
