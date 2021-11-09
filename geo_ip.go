package main

import (
	"context"
	_ "embed"
	"net"

	"github.com/jackc/pgx/v4"
	"github.com/oschwald/geoip2-golang"
)

// GeoResult is result of ip and client geo data
type GeoResult struct {
	// ip data
	IP                       string
	IPCountry                string
	IPCity                   string
	IPCityGeoID              uint
	IPLocationLatitude       float64
	IPLocationLongitude      float64
	IPLocationAccuracyRadius uint16
	IPASNumber               uint
	IPASName                 string

	// geo location data
	ClientGeoLatitude       float64
	ClientGeoLongitude      float64
	ClientGeoAccuracyRadius float64
	ClientCountry           string
	ClientCity              string
	ClientAdministratorArea string

	// if client geo other wise use ip address
	Country           string
	AdministratorArea string
	City              string
}

// GeoParser is instace for parse ip data
type GeoParser struct {
	postgisConnection *pgx.Conn
	geoASN            *geoip2.Reader
	geoCity           *geoip2.Reader
}

// GeoIPResult is result of geo location base on ip
type GeoIPResult struct {
	// base on client IP
	IP                       string
	IPCountry                string
	IPCity                   string
	IPCityGeoID              uint
	IPLocationLatitude       float64
	IPLocationLongitude      float64
	IPLocationAccuracyRadius uint16
	IPASNumber               uint
	IPASName                 string
}

// GeoLocationResult is result of geo location base on client geo data
type GeoLocationResult struct {
	// base on client geo location
	ClientGeoLatitude       float64
	ClientGeoLongitude      float64
	ClientGeoAccuracyRadius float64
	ClientCountry           string
	ClientCity              string
	ClientAdministratorArea string
}

// NewGeoIPParser will return new instance of GeoIPParser
func NewGeoIPParser(
	postgisConnection *pgx.Conn,
	mmdbCityPath string,
	mmdbASNPath string,
) *GeoParser {
	dbCity, err := geoip2.Open(mmdbCityPath)
	if err != nil {
		panic(err)
	}
	dbASN, err := geoip2.Open(mmdbASNPath)
	if err != nil {
		panic(err)
	}

	geoIPParser := GeoParser{
		postgisConnection: postgisConnection,
		geoASN:            dbASN,
		geoCity:           dbCity,
	}

	return &geoIPParser
}

// Parse parsing geo data from ip
func (geoParser *GeoParser) ParseIP(ipString string) *GeoIPResult {
	result := GeoIPResult{}

	ip := net.ParseIP(ipString)

	if ip == nil {
		return &result
	}

	result.IP = ipString

	recordCity, err := geoParser.geoCity.City(ip)
	if err == nil {
		result.IPCountry = recordCity.Country.IsoCode
		result.IPCity = recordCity.City.Names["en"]
		result.IPCityGeoID = recordCity.City.GeoNameID
		result.IPLocationLongitude = recordCity.Location.Longitude
		result.IPLocationLatitude = recordCity.Location.Latitude
		result.IPLocationAccuracyRadius = recordCity.Location.AccuracyRadius
	}

	recordASN, err := geoParser.geoASN.ASN(ip)
	if err == nil {
		result.IPASName = recordASN.AutonomousSystemOrganization
		result.IPASNumber = recordASN.AutonomousSystemNumber
	}

	return &result
}

// Parse parsing geo data from ip
func (geoParser *GeoParser) LookupLocation(
	clientGeoLatitude float64,
	clientGeoLongitude float64,
	clientGeoAccuracyRadius float64,
) *GeoLocationResult {
	result := GeoLocationResult{
		ClientGeoLatitude:       clientGeoLatitude,
		ClientGeoLongitude:      clientGeoLongitude,
		ClientGeoAccuracyRadius: clientGeoAccuracyRadius,
	}

	query := `
		SELECT
			"geo"."name",
			"adminCode"."name",
			"countryInfo"."iso"
		FROM "geo"
		LEFT JOIN "countryInfo" ON ("countryInfo"."geonameid" = "geo"."country")
		LEFT JOIN "adminCode" ON ("adminCode"."id" = "geo"."adminCode")
		ORDER BY "geo"."location" <-> ST_SetSRID(ST_MakePoint($1, $2), 4326) LIMIT 1;
	`

	geoParser.postgisConnection.QueryRow(context.Background(), query, clientGeoLongitude, clientGeoLatitude).Scan(
		&result.ClientCity,
		&result.ClientAdministratorArea,
		&result.ClientCountry,
	)

	return &result
}

// NewGeoResult
func NewGeoResult(
	geoIPResult *GeoIPResult,
) *GeoResult {
	result := GeoResult{
		IP:                       geoIPResult.IP,
		IPCountry:                geoIPResult.IPCountry,
		IPCity:                   geoIPResult.IPCity,
		IPCityGeoID:              geoIPResult.IPCityGeoID,
		IPLocationLatitude:       geoIPResult.IPLocationLatitude,
		IPLocationLongitude:      geoIPResult.IPLocationLongitude,
		IPLocationAccuracyRadius: geoIPResult.IPLocationAccuracyRadius,
		IPASNumber:               geoIPResult.IPASNumber,
		IPASName:                 geoIPResult.IPASName,
		Country:                  geoIPResult.IPCountry,
		City:                     geoIPResult.IPCity,
	}

	return &result
}

func (geoResult *GeoResult) AddLocation(geoLocationResult *GeoLocationResult) {
	// AdministratorArea
	geoResult.AdministratorArea = geoLocationResult.ClientAdministratorArea
	// Country
	geoResult.ClientCountry = geoLocationResult.ClientCountry
	geoResult.Country = geoLocationResult.ClientCountry
	// City
	geoResult.ClientCity = geoLocationResult.ClientCity
	geoResult.City = geoLocationResult.ClientCity
}
