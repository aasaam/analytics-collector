package main

import (
	"net/url"
	"testing"
)

func TestParseUTM1(t *testing.T) {
	u := getURL("https://www.google.com")

	utm := parseUTM(u)
	if utm.UTMValid || utm.UTMExist {
		t.Errorf("invalid utm parse")
	}
}
func TestParseUTM2(t *testing.T) {
	u := getURL("https://www.example.com?UTM_SOURCE=source&utm_medium=medium&utm_campaign=sale1&utm_id=id&utm_term=keyword1&utm_content=content")

	utm := parseUTM(u)
	if !utm.UTMValid || !utm.UTMExist {
		t.Errorf("invalid utm parse")
	}
}
func TestParseUTM3(t *testing.T) {
	sampleName := " This IS <>\"' for Name"
	u := getURL("https://www.example.com?UTM_SOURCE=source&utm_medium=medium&utm_campaign=" + sampleName + "&utm_id=id&utm_term=keyword1&utm_content=content")

	utm := parseUTM(u)
	if !utm.UTMValid || !utm.UTMExist {
		t.Errorf("invalid utm parse")
	}
}
func TestParseUTM4(t *testing.T) {
	sampleName := " This IS <>\"' for Name"
	u := getURL("https://www.example.com/?utm_source=source&utm_medium=medium&utm_campaign=" + url.QueryEscape(sampleName) + "&utm_id=id&utm_term=keyword1&utm_content=content")

	utm := parseUTM(u)
	if !utm.UTMValid || !utm.UTMExist {
		t.Errorf("invalid utm parse")
	}
}
func TestParseUTM5(t *testing.T) {
	u := getURL("https://www.example.com?UTM_SOURCE=source&utm_medium=medium")

	utm := parseUTM(u)
	if utm.UTMValid || !utm.UTMExist {
		t.Errorf("invalid utm parse")
	}
}
func TestParseUTM6(t *testing.T) {
	utm := parseUTM(nil)
	if utm.UTMValid || utm.UTMExist {
		t.Errorf("invalid utm parse")
	}
}

func BenchmarkParseUTM(b *testing.B) {
	u := getURL("https://www.example.com?UTM_SOURCE=source&utm_medium=medium&utm_campaign=sale1&utm_id=id&utm_term=keyword1&utm_content=content")
	for n := 0; n < b.N; n++ {
		parseUTM(u)
	}
}
