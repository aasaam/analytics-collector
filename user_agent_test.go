package main

import (
	"testing"
)

func TestUserAgentParser1(t *testing.T) {
	uaStrings := []string{
		"Mozilla/5.0 (Linux; Android 5.0.2; SAMSUNG SM-A500FU Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/3.3 Chrome/38.0.2125.102 Mobile Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		"curl/7.68.0",
		"AdsBot-Google (+http://www.google.com/adsbot.html)",
		"Wget/1.20.3 (linux-gnu)",
	}

	for _, uaString := range uaStrings {
		uaParser := NewUserAgentParser()

		parsed := uaParser.Parse(uaString)
		prettyPrint(parsed)
		if parsed.BrowserName == "" {
			t.Errorf("invalid parsed data")
		}
	}
}

func BenchmarkNewUserAgentParser(b *testing.B) {
	uaString := "Mozilla/5.0 (Linux; Android 5.0.2; SAMSUNG SM-A500FU Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/3.3 Chrome/38.0.2125.102 Mobile Safari/537.36"
	uaParser := NewUserAgentParser()
	for n := 0; n < b.N; n++ {
		uaParser.Parse(uaString)
	}
}
