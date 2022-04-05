package main

import (
	"net/url"
	"testing"
)

func TestParseUtm1(t *testing.T) {
	u := getURL("https://www.google.com")

	utm := parseUTM(u)
	if utm.UtmValid || utm.UtmExist {
		t.Errorf("invalid utm parse")
	}
}
func TestParseUtm2(t *testing.T) {
	u := getURL("https://www.example.com?Utm_SOURCE=source&utm_medium=medium&utm_campaign=sale1&utm_id=id&utm_term=keyword1&utm_content=content")

	utm := parseUTM(u)
	if !utm.UtmValid || !utm.UtmExist {
		t.Errorf("invalid utm parse")
	}
}
func TestParseUtm3(t *testing.T) {
	sampleName := " This IS <>\"' for Name"
	u := getURL("https://www.example.com?Utm_SOURCE=source&utm_medium=medium&utm_campaign=" + sampleName + "&utm_id=id&utm_term=keyword1&utm_content=content")

	utm := parseUTM(u)
	if !utm.UtmValid || !utm.UtmExist {
		t.Errorf("invalid utm parse")
	}
}
func TestParseUtm4(t *testing.T) {
	sampleName := " This IS <>\"' for Name"
	u := getURL("https://www.example.com/?utm_source=source&utm_medium=medium&utm_campaign=" + url.QueryEscape(sampleName) + "&utm_id=id&utm_term=keyword1&utm_content=content")

	utm := parseUTM(u)
	if !utm.UtmValid || !utm.UtmExist {
		t.Errorf("invalid utm parse")
	}
}
func TestParseUtm5(t *testing.T) {
	u := getURL("https://www.example.com?Utm_SOURCE=source&utm_medium=medium")

	utm := parseUTM(u)
	if utm.UtmValid || !utm.UtmExist {
		t.Errorf("invalid utm parse")
	}
}
func TestParseUtm6(t *testing.T) {
	utm := parseUTM(nil)
	if utm.UtmValid || utm.UtmExist {
		t.Errorf("invalid utm parse")
	}
}

func BenchmarkParseUtm(b *testing.B) {
	u := getURL("https://www.example.com?Utm_SOURCE=source&utm_medium=medium&utm_campaign=sale1&utm_id=id&utm_term=keyword1&utm_content=content")
	for n := 0; n < b.N; n++ {
		parseUTM(u)
	}
}
