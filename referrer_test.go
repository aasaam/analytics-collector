package main

import (
	"testing"
)

func TestReferrerParser1(t *testing.T) {
	parser := NewReferrerParser()

	data := parser.Parse("https://www.google.com", "https://www.google.com")
	if data.ExternalHost != "" {
		t.Errorf("invalid data")
	}

	prettyPrint(data)

	if data.ReferrerType != ReferrerTypeSearchEngine {
		t.Errorf("invalid data")
	}

	data2 := parser.Parse("http://www.example.com/page.html", "https://www.google.com")
	prettyPrint(data2)
	if data2.ExternalHost != "www.google.com" {
		t.Errorf("invalid data")
	}
}
func TestReferrerParser2(t *testing.T) {
	parser := NewReferrerParser()

	data := parser.Parse("https://www.google.com", "android-app://com.google.android.gm/")
	if data.ExternalHost == "" {
		t.Errorf("invalid data")
	}

	prettyPrint(data)

	if data.ReferrerType != ReferrerTypeMailProvider {
		t.Errorf("invalid data")
	}
}
func TestReferrerParser3(t *testing.T) {
	parser := NewReferrerParser()

	data := parser.Parse("1", "2")
	if data.ExternalHost != "" {
		t.Errorf("invalid data")
	}
}

func BenchmarkReferrerParser(b *testing.B) {
	parser := NewReferrerParser()
	for n := 0; n < b.N; n++ {
		parser.Parse("https://www.google.com", "https://www.google.com")
	}
}
