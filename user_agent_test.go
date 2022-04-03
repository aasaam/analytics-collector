package main

import (
	_ "embed"
	"testing"
)

func TestUserAgentParser1(t *testing.T) {

	uaStrings := []string{
		"Mozilla/5.0 (Linux; Android 5.0.2; SAMSUNG SM-A500FU Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/3.3 Chrome/38.0.2125.102 Mobile Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		"Mozilla/5.0 (iPad; CPU OS 7_0 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) CriOS/30.0.1599.12 Mobile/11A465 Safari/8536.25 (3B92C18B-D9DE-4CB7-A02A-22FD2AF17C8F)",
		"curl/7.68.0",
		"AdsBot-Google (+http://www.google.com/adsbot.html)",
		"Wget/1.20.3 (linux-gnu)",
	}

	uaParser := newUserAgentParser()

	for _, uaString := range uaStrings {

		parsed := uaParser.parse(uaString)
		if parsed.UaBrowserName == "" {
			t.Errorf("invalid parsed data")
		}
	}
}

func BenchmarkNewUserAgentParser(b *testing.B) {
	uaString := "Mozilla/5.0 (Linux; Android 5.0.2; SAMSUNG SM-A500FU Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/3.3 Chrome/38.0.2125.102 Mobile Safari/537.36"
	uaParser := newUserAgentParser()
	for n := 0; n < b.N; n++ {
		uaParser.parse(uaString)
	}
}
