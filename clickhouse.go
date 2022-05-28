package main

import (
	"context"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/paulmach/orb"
)

func getClickhouseConnection(
	servers string,
	database string,
	username string,
	password string,
	maxExecutionTime int,
	dialTimeout int,
	debug bool,
	compressionLZ4 bool,
	maxIdleConns int,
	maxOpenConns int,
	connMaxLifetime int,
	maxBlockSize int,
	progress func(p *clickhouse.Progress),
	profile func(p *clickhouse.ProfileInfo),
) (driver.Conn, context.Context, error) {

	clickhouseOpts := clickhouse.Options{
		Addr: strings.Split(servers, ","),
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": maxExecutionTime,
		},
		DialTimeout:     time.Duration(dialTimeout) * time.Second,
		Debug:           debug,
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: time.Duration(connMaxLifetime) * time.Second,
	}

	if compressionLZ4 {
		clickhouseOpts.Compression = &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		}
	}

	clickhouseConn, clickhouseConnErr := clickhouse.Open(&clickhouseOpts)

	if clickhouseConnErr != nil {
		return nil, nil, clickhouseConnErr
	}

	ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": maxBlockSize,
	}), clickhouse.WithProgress(progress), clickhouse.WithProfileInfo(profile))

	if pingErr := clickhouseConn.Ping(ctx); pingErr != nil {
		return nil, nil, pingErr
	}

	return clickhouseConn, ctx, nil
}

func insertClientErrBatch(
	batch driver.Batch,
	rec record,
) error {
	geoIPPoint := orb.Point{rec.GeoResult.GeoIPLocationLongitude, rec.GeoResult.GeoIPLocationLatitude}

	return batch.Append(
		rec.ClientErrorMessage,
		rec.ClientErrorObject,

		getURLString(rec.PURL),

		// geo:asn
		rec.GeoResult.GeoIPAutonomousSystemNumber,
		rec.GeoResult.GeoIPAutonomousSystemOrganization,
		// geo:ip
		rec.GeoResult.GeoIPAdministratorArea,
		rec.GeoResult.GeoIPCity,
		rec.GeoResult.GeoIPCityGeoNameID,
		rec.GeoResult.GeoIPCountry,
		rec.GeoResult.GeoIPLocationLatitude,
		rec.GeoResult.GeoIPLocationLongitude,
		// geo:ip+
		geoIPPoint,

		// user agent
		rec.UserAgentResult.UaType,
		rec.UserAgentResult.UaFull,
		rec.UserAgentResult.UaChecksum,
		rec.UserAgentResult.UaBrowserName,
		rec.UserAgentResult.UaBrowserVersionMajor,
		rec.UserAgentResult.UaBrowserVersion,
		rec.UserAgentResult.UaOSName,
		rec.UserAgentResult.UaOSVersionMajor,
		rec.UserAgentResult.UaOSVersion,
		rec.UserAgentResult.UaDeviceBrand,
		rec.UserAgentResult.UaDeviceFamily,
		rec.UserAgentResult.UaDeviceModel,

		rec.IP,
		rec.PublicInstanceID,
		rec.Mode,
		rec.Created,
	)
}
func insertRecordBatch(
	batch driver.Batch,
	rec record,
	ECategory string,
	EAction string,
	ELabel string,
	EIdent string,
	EValue uint64,
) error {
	geoIPPoint := orb.Point{rec.GeoResult.GeoIPLocationLongitude, rec.GeoResult.GeoIPLocationLatitude}
	geoClientPoint := orb.Point{rec.GeoResult.GeoClientLocationLongitude, rec.GeoResult.GeoClientLocationLatitude}
	geoResultPoint := orb.Point{rec.GeoResult.GeoResultLocationLongitude, rec.GeoResult.GeoResultLocationLatitude}

	return batch.Append(

		// event
		ECategory,
		EAction,
		ELabel,
		EIdent,
		EValue,

		// custom segments
		rec.Segments.S1N,
		rec.Segments.S2N,
		rec.Segments.S3N,
		rec.Segments.S4N,
		rec.Segments.S5N,
		rec.Segments.S1V,
		rec.Segments.S2V,
		rec.Segments.S3V,
		rec.Segments.S4V,
		rec.Segments.S5V,

		// page
		boolUint8(rec.PIsIframe),
		boolUint8(rec.PIsTouchSupport),
		getURLString(rec.PURL),
		checksum(getURLString(rec.PURL)),
		rec.PTitle,
		getURLString(rec.PCanonicalURL),
		checksum(getURLString(rec.PCanonicalURL)),
		rec.PLang,
		rec.PEntityID,
		rec.PEntityModule,
		rec.PEntityTaxonomyID,
		rec.PKeywords,

		// referer
		getURLString(rec.PRefererURL.RefURL),
		rec.PRefererURL.RefExternalHost,
		rec.PRefererURL.RefExternalDomain,
		rec.PRefererURL.RefName,
		rec.PRefererURL.RefProtocol,
		rec.PRefererURL.RefType,

		// session referer
		getURLString(rec.SRefererURL.RefURL),
		rec.SRefererURL.RefExternalHost,
		rec.SRefererURL.RefExternalDomain,
		rec.SRefererURL.RefName,
		rec.SRefererURL.RefProtocol,
		rec.SRefererURL.RefType,

		// utm
		boolUint8(rec.Utm.UtmValid),
		boolUint8(rec.Utm.UtmExist),
		rec.Utm.UtmSource,
		rec.Utm.UtmMedium,
		rec.Utm.UtmCampaign,
		rec.Utm.UtmID,
		rec.Utm.UtmTerm,
		rec.Utm.UtmContent,

		// performance
		boolUint8(rec.Performance.PerfIsProcessed),
		rec.Performance.PerfPageLoadTime,
		rec.Performance.PerfDomainLookupTime,
		rec.Performance.PerfTCPConnectTime,
		rec.Performance.PerfServerResponseTime,
		rec.Performance.PerfPageDownloadTime,
		rec.Performance.PerfRedirectTime,
		rec.Performance.PerfDOMInteractiveTime,
		rec.Performance.PerfContentLoadTime,
		rec.Performance.PerfResource,

		// breadcrumb
		rec.BreadCrumb.BCLevel,
		rec.BreadCrumb.BCN1,
		rec.BreadCrumb.BCN2,
		rec.BreadCrumb.BCN3,
		rec.BreadCrumb.BCN4,
		rec.BreadCrumb.BCN5,
		rec.BreadCrumb.BCP1,
		rec.BreadCrumb.BCP2,
		rec.BreadCrumb.BCP3,
		rec.BreadCrumb.BCP4,
		rec.BreadCrumb.BCP5,

		// user agent
		rec.UserAgentResult.UaType,
		rec.UserAgentResult.UaFull,
		rec.UserAgentResult.UaChecksum,
		rec.UserAgentResult.UaBrowserName,
		rec.UserAgentResult.UaBrowserVersionMajor,
		rec.UserAgentResult.UaBrowserVersion,
		rec.UserAgentResult.UaOSName,
		rec.UserAgentResult.UaOSVersionMajor,
		rec.UserAgentResult.UaOSVersion,
		rec.UserAgentResult.UaDeviceBrand,
		rec.UserAgentResult.UaDeviceFamily,
		rec.UserAgentResult.UaDeviceModel,

		// screen
		boolUint8(rec.ScreenInfo.ScrScreenOrientation),
		boolUint8(rec.ScreenInfo.ScrScreenOrientationIsPortrait),
		boolUint8(rec.ScreenInfo.ScrScreenOrientationIsSecondary),
		rec.ScreenInfo.ScrScreen,
		rec.ScreenInfo.ScrScreenWidth,
		rec.ScreenInfo.ScrScreenHeight,
		rec.ScreenInfo.ScrViewport,
		rec.ScreenInfo.ScrViewportWidth,
		rec.ScreenInfo.ScrViewportHeight,
		rec.ScreenInfo.ScrResoluton,
		rec.ScreenInfo.ScrResolutonWidth,
		rec.ScreenInfo.ScrResolutonHeight,
		rec.ScreenInfo.ScrDevicePixelRatio,
		rec.ScreenInfo.ScrColorDepth,

		// geo:asn
		rec.GeoResult.GeoIPAutonomousSystemNumber,
		rec.GeoResult.GeoIPAutonomousSystemOrganization,
		// geo:ip
		rec.GeoResult.GeoIPAdministratorArea,
		rec.GeoResult.GeoIPCity,
		rec.GeoResult.GeoIPCityGeoNameID,
		rec.GeoResult.GeoIPCountry,
		rec.GeoResult.GeoIPLocationLatitude,
		rec.GeoResult.GeoIPLocationLongitude,
		// geo:ip+
		geoIPPoint,
		// geo:client
		rec.GeoResult.GeoClientAdministratorArea,
		rec.GeoResult.GeoClientCity,
		rec.GeoResult.GeoClientCityGeoNameID,
		rec.GeoResult.GeoClientCountry,
		rec.GeoResult.GeoClientLocationLatitude,
		rec.GeoResult.GeoClientLocationLongitude,
		// geo:client+
		geoClientPoint,
		// geo:result
		boolUint8(rec.GeoResult.GeoResultFromClient),
		rec.GeoResult.GeoResultAdministratorArea,
		rec.GeoResult.GeoResultCity,
		rec.GeoResult.GeoResultCityGeoNameID,
		rec.GeoResult.GeoResultCountry,
		rec.GeoResult.GeoResultLocationLatitude,
		rec.GeoResult.GeoResultLocationLongitude,
		// geo:result+
		geoResultPoint,

		rec.CID.CidType,
		rec.CID.CidUserChecksum,
		rec.CID.CidSessionChecksum,
		rec.CID.CidStdInitTime,
		rec.CID.CidStdSessionTime,

		rec.IP,
		rec.PublicInstanceID,
		rec.Mode,
		rec.CursorID,
		rec.Created,
	)
}
