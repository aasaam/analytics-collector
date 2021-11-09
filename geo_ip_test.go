package main

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
)

var sampleIP1 = "81.2.69.142"
var sampleIP2 = "12.81.92.2"

func TestGeoIPParser1(t *testing.T) {
	postgisConnectionString := os.Getenv("POSTGIS_URI")
	if postgisConnectionString == "" {
		postgisConnectionString = "postgres://geonames:geonames@127.0.0.1:5432/geonames"
	}

	conn, err := pgx.Connect(context.Background(), postgisConnectionString)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	geoIPParser := NewGeoIPParser(conn, "data/Geo-City.mmdb", "data/Geo-ASN.mmdb")
	parsedIP1 := geoIPParser.ParseIP(sampleIP1)

	if parsedIP1.IPCountry != "GB" {
		t.Errorf("invalid parsed data")
	}

	parsedIP2 := geoIPParser.ParseIP(sampleIP2)
	if parsedIP2.IPASName != "AT&T Services" {
		t.Errorf("invalid parsed data")
	}

	parsedIP3 := geoIPParser.ParseIP("z.z.z.z")
	if parsedIP3.IP != "" {
		t.Errorf("invalid parsed data")
	}

	lookup1 := geoIPParser.LookupLocation(35.6892, 51.3890, 1000)
	if lookup1.ClientCity != "Tehran" {
		t.Errorf("invalid parsed data")
	}

	// default location data base on IP address
	geoResult := NewGeoResult(parsedIP1)
	if geoResult.IPCountry != "GB" || geoResult.Country != "GB" {
		t.Errorf("invalid parsed data")
	}

	// now patch with client geo location
	geoResult.AddLocation(lookup1)
	if geoResult.ClientCity != "Tehran" || geoResult.Country != "IR" {
		t.Errorf("invalid parsed data")
	}
}
