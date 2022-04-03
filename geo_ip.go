package main

import (
	"context"
	_ "embed"
	"errors"
	"net"

	"github.com/jackc/pgx/v4"
	"github.com/oschwald/geoip2-golang"
)

type geoParser struct {
	postgisConnection *pgx.Conn
	geoASN            *geoip2.Reader
	geoCity           *geoip2.Reader
}

type geonameData struct {
	valid             bool
	id                uint
	city              string
	country           string
	administratorArea string
}

type geoResult struct {
	GeoIsProcessed bool
	GeoIP          string

	// geo:asn
	GeoIPAutonomousSystemNumber       uint16
	GeoIPAutonomousSystemOrganization string

	// geo:ip
	GeoIPAdministratorArea string
	GeoIPCity              string
	GeoIPCityGeoNameID     uint16
	GeoIPCountry           string
	GeoIPLocationLatitude  float64
	GeoIPLocationLongitude float64

	// geo:client
	GeoClientAdministratorArea string
	GeoClientCity              string
	GeoClientCityGeoNameID     uint16
	GeoClientCountry           string
	GeoClientLocationLatitude  float64
	GeoClientLocationLongitude float64

	GeoResultAdministratorArea string
	GeoResultCity              string
	GeoResultCityGeoNameID     uint16
	GeoResultCountry           string
	GeoResultFromClient        bool
	GeoResultLocationLatitude  float64
	GeoResultLocationLongitude float64
}

func newGeoParser(
	postgisConnection *pgx.Conn,
	mmdbCityPath string,
	mmdbASNPath string,
) (*geoParser, error) {
	dbCity, err := geoip2.Open(mmdbCityPath)
	if err != nil {
		return nil, err
	}
	dbASN, err := geoip2.Open(mmdbASNPath)
	if err != nil {
		return nil, err
	}

	geoIPParser := geoParser{
		postgisConnection: postgisConnection,
		geoASN:            dbASN,
		geoCity:           dbCity,
	}

	return &geoIPParser, nil
}

func (geoParser *geoParser) getGeonameData(
	latitude float64,
	longitude float64,
) (geonameData, error) {
	query := `
		SELECT
			"geo"."geonameid" AS "id",
			"geo"."name" AS "city",
			"adminCode"."name" AS "admin",
			"countryInfo"."iso" AS "country"
		FROM "geo"
		LEFT JOIN "countryInfo" ON ("countryInfo"."geonameid" = "geo"."country")
		LEFT JOIN "adminCode" ON ("adminCode"."id" = "geo"."adminCode")
		ORDER BY "geo"."location" <-> ST_SetSRID(ST_MakePoint($1, $2), 4326) LIMIT 1;
	`

	result := geonameData{}

	err := geoParser.postgisConnection.QueryRow(context.Background(), query, longitude, latitude).Scan(
		&result.id,
		&result.city,
		&result.administratorArea,
		&result.country,
	)

	if err == nil && result.city != "" {
		result.valid = true
		return result, nil
	}

	return geonameData{}, errors.New("geoname data not found")
}

func (geoParser *geoParser) newResultFromIP(ip net.IP) geoResult {
	obj := geoResult{}

	if ip == nil || ip.To4() == nil {
		return obj
	}

	obj.GeoIP = ip.String()
	obj.GeoIsProcessed = true
	obj.GeoResultFromClient = false

	recordCity, err := geoParser.geoCity.City(ip)

	if err == nil {
		obj.GeoIPCity = recordCity.City.Names["en"]
		obj.GeoIPCityGeoNameID = uint16(recordCity.City.GeoNameID)
		obj.GeoIPCountry = recordCity.Country.IsoCode
		obj.GeoIPLocationLatitude = recordCity.Location.Latitude
		obj.GeoIPLocationLongitude = recordCity.Location.Longitude

		var geonameData geonameData
		if obj.GeoIPLocationLongitude != 0 || obj.GeoIPLocationLatitude != 0 {
			geonameData, _ = geoParser.getGeonameData(obj.GeoIPLocationLatitude, obj.GeoIPLocationLongitude)
		}

		if geonameData.valid {
			obj.GeoIPAdministratorArea = geonameData.administratorArea
			if obj.GeoIPCity == "" {
				obj.GeoIPCity = geonameData.city
			}
			if obj.GeoIPCityGeoNameID < 1 {
				obj.GeoIPCityGeoNameID = uint16(geonameData.id)
			}
		}

		obj.GeoResultAdministratorArea = obj.GeoIPAdministratorArea
		obj.GeoResultCity = obj.GeoIPCity
		obj.GeoResultCityGeoNameID = obj.GeoIPCityGeoNameID
		obj.GeoResultCountry = obj.GeoIPCountry
		obj.GeoResultLocationLatitude = obj.GeoIPLocationLatitude
		obj.GeoResultLocationLongitude = obj.GeoIPLocationLongitude
	}

	recordASN, err := geoParser.geoASN.ASN(ip)
	if err == nil {
		obj.GeoIPAutonomousSystemOrganization = recordASN.AutonomousSystemOrganization
		obj.GeoIPAutonomousSystemNumber = uint16(recordASN.AutonomousSystemNumber)
	}

	return obj
}

func (geoParser *geoParser) clientLocationUpdate(
	obj geoResult,
	clientLocationLatitude float64,
	clientLocationLongitude float64,
) geoResult {

	query := `
		SELECT
			"geo"."geonameid" AS "id",
			"geo"."name" AS "city",
			"adminCode"."name" AS "admin",
			"countryInfo"."iso" AS "country"
		FROM "geo"
		LEFT JOIN "countryInfo" ON ("countryInfo"."geonameid" = "geo"."country")
		LEFT JOIN "adminCode" ON ("adminCode"."id" = "geo"."adminCode")
		ORDER BY "geo"."location" <-> ST_SetSRID(ST_MakePoint($1, $2), 4326) LIMIT 1;
	`

	result := geoResult{
		GeoClientLocationLatitude:  clientLocationLatitude,
		GeoClientLocationLongitude: clientLocationLongitude,
	}

	err := geoParser.postgisConnection.QueryRow(context.Background(), query, clientLocationLongitude, clientLocationLatitude).Scan(
		&result.GeoClientCityGeoNameID,
		&result.GeoClientCity,
		&result.GeoClientAdministratorArea,
		&result.GeoClientCountry,
	)

	if err == nil && result.GeoClientCountry != "" {
		obj.GeoResultFromClient = true
		obj.GeoClientAdministratorArea = result.GeoClientAdministratorArea
		obj.GeoClientCity = result.GeoClientCity
		obj.GeoClientCityGeoNameID = result.GeoClientCityGeoNameID
		obj.GeoClientCountry = result.GeoClientCountry
		obj.GeoClientLocationLatitude = result.GeoClientLocationLatitude
		obj.GeoClientLocationLongitude = result.GeoClientLocationLongitude

		obj.GeoResultAdministratorArea = obj.GeoClientAdministratorArea
		obj.GeoResultCity = obj.GeoClientCity
		obj.GeoResultCityGeoNameID = obj.GeoClientCityGeoNameID
		obj.GeoResultCountry = obj.GeoClientCountry
		obj.GeoResultLocationLatitude = obj.GeoClientLocationLatitude
		obj.GeoResultLocationLongitude = obj.GeoClientLocationLongitude
	}

	return obj
}
