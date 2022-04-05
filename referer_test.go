package main

import (
	"testing"
)

func TestRefererParser1(t *testing.T) {
	refererParser := newRefererParser()

	data1 := refererParser.parse(getURL("https://www.google.com"), getURL("https://www.google.com"))
	if data1.RefExternalHost != "" {
		t.Errorf("invalid data")
	}
	if data1.RefType == refererTypeSearchEngine {
		t.Errorf("invalid data")
	}

	data2 := refererParser.parse(getURL("http://www.example.com/page.html"), getURL("https://www.google.com"))
	if data2.RefExternalHost != "www.google.com" {
		t.Errorf("invalid data")
	}
}

func TestRefererParser2(t *testing.T) {
	refererParser := newRefererParser()

	data := refererParser.parse(getURL("https://www.google.com"), getURL("android-app://com.google.android.gm/"))
	if data.RefExternalHost == "" {
		t.Errorf("invalid data")
	}

	if data.RefType != refererTypeMailProvider {
		t.Errorf("invalid data")
	}
}
func TestRefererParser3(t *testing.T) {
	refererParser := newRefererParser()

	data := refererParser.parse(getURL("1"), getURL("2"))
	if data.RefExternalHost != "" {
		t.Errorf("invalid data")
	}
}

func TestRefererParser4(t *testing.T) {
	refererParser := newRefererParser()

	data := refererParser.parse(getURL("\x18"), getURL("https://www.google.com"))
	if data.RefExternalHost != "" {
		t.Errorf("invalid data")
	}
}

func TestRefererParser5(t *testing.T) {
	refererParser := newRefererParser()

	data := refererParser.parse(getURL("https://www.google.com"), getURL("\x18"))
	if data.RefExternalHost != "" {
		t.Errorf("invalid data")
	}
}

func TestRefererParser6(t *testing.T) {
	refererParser := newRefererParser()

	data1 := refererParser.parse(getURL("https://www.example.com/path/foo/bar"), getURL("https://www.search-engine.com/path/foo/bar"))
	if data1.RefExternalHost == "" || data1.RefExternalDomain == "" {
		t.Errorf("invalid data")
	}

	data2 := refererParser.parse(getURL("https://www.example.com/path/foo/bar"), getURL("https://another-sub.example.com/path/foo/bar"))
	if data2.RefExternalHost == "" || data2.RefExternalDomain != "" {
		t.Errorf("invalid data")
	}

	data3 := refererParser.parse(getURL("https://www.example.com/path/foo/bar"), getURL("https://www.example.com"))
	if data3.RefExternalHost != "" || data3.RefExternalDomain != "" {
		t.Errorf("invalid data")
	}
}

func BenchmarkRefererParser(b *testing.B) {
	refererParser := newRefererParser()
	u1 := getURL("https://www.google.com")
	u2 := getURL("https://www.google.com")
	for n := 0; n < b.N; n++ {
		refererParser.parse(u1, u2)
	}
}
