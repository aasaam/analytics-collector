package main

import (
	"context"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4"
)

var sampleIPString = "1.1.1.1"

func getGeoParser() *geoParser {
	postgisConnectionString := os.Getenv("POSTGIS_URI")
	if postgisConnectionString == "" {
		postgisConnectionString = "postgres://geonames:geonames@127.0.0.1:5432/geonames"
	}

	conn, err := pgx.Connect(context.Background(), postgisConnectionString)
	if err != nil {
		panic(err)
	}
	geoParser, err := newGeoParser(conn, "tmp/GeoLite2-City.mmdb", "tmp/GeoLite2-ASN.mmdb")
	if err != nil {
		panic(err)
	}
	return geoParser
}

func TestGeoIPParser1(t *testing.T) {
	geoParser := getGeoParser()

	ip := net.ParseIP(sampleIPString)

	parsedIP1 := geoParser.newResultFromIP(ip)
	if len(parsedIP1.GeoIPCountry) != 2 {
		t.Errorf("invalid parsed data")
	}

	parsedIP2 := geoParser.newResultFromIP(ip)
	if !strings.Contains(strings.ToLower(parsedIP2.GeoIPAutonomousSystemOrganization), "cloudflare") {
		t.Errorf("invalid parsed data")
	}

	parsedIP3 := geoParser.newResultFromIP(nil)
	if parsedIP3.GeoIP != "" {
		t.Errorf("invalid parsed data")
	}

	parsedIP4 := geoParser.newResultFromIP(net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"))
	if parsedIP4.GeoIP != "" {
		t.Errorf("invalid parsed data")
	}
}

func TestGeoIPParser2(t *testing.T) {
	geoParser := getGeoParser()

	ip := net.ParseIP(sampleIPString)

	// ip from England
	parsedIP1 := geoParser.newResultFromIP(ip)
	if len(parsedIP1.GeoIPCountry) != 2 {
		t.Errorf("invalid parsed data")
	}

	// client from iran tehran
	parsedIP1 = geoParser.clientLocationUpdate(parsedIP1, 35.6892, 51.3890)

	if parsedIP1.GeoResultCountry != "IR" {
		t.Errorf("invalid parsed data")
	}
}
