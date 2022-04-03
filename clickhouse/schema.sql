SET allow_experimental_geo_types = 1;
CREATE DATABASE IF NOT EXISTS analytics;

CREATE TABLE IF NOT EXISTS analytics.ClientErrors
(
  Msg                                         String,
  Err                                         String,

  -- requirements
  IP                                          IPv4,
	PublicInstaceID                             String,
	Mode                                        UInt8,
	Created                                     Datetime
)
ENGINE = MergeTree()
ORDER BY (Created, xxHash32(Msg))
PARTITION BY toYYYYMM(Created)
SAMPLE BY (xxHash32(msg));

CREATE TABLE IF NOT EXISTS analytics.Records
(
  -- etc
  PropsUserIDOrName                           String,

  -- page
  PIsIframe                                   UInt8, --bool
  PIsTouchSupport                             UInt8, --bool
  PURL                                        String,
  PTitle                                      String,
  PCanonicalURL                               String,
  PLang                                       LowCardinality(String),
  PEntityID                                   String,
  PEntityModule                               String,
  PEntityTaxonomyID                           LowCardinality(String),
  PKeywords                                   Array(String),

  -- referer
  PRefererURLURL                              String,
  PRefererURLExternalHost                     String,
  PRefererURLExternalName                     String,
  PRefererURLExternalType                     UInt8,

  -- session referer
  SRefererURLURL                              String,
  SRefererURLExternalHost                     String,
  SRefererURLExternalName                     String,
  SRefererURLExternalType                     UInt8,

  -- performance
  PerfExist                                   UInt8, --bool
	PerfPageLoadTime                            UInt16,
	PerfDomainLookupTime                        UInt16,
	PerfTCPConnectTime                          UInt16,
	PerfServerResponseTime                      UInt16,
	PerfPageDownloadTime                        UInt16,
	PerfRedirectTime                            UInt16,
	PerfDOMInteractiveTime                      UInt16,
	PerfContentLoadTime                         UInt16,
	PerfResource                                UInt16,

  -- breadcrumb
  BCN1                                        String,
	BCN2                                        String,
	BCN3                                        String,
	BCN4                                        String,
	BCN5                                        String,
	BCP1                                        String,
	BCP2                                        String,
	BCP3                                        String,
	BCP4                                        String,
	BCP5                                        String,

  -- event
  ECategory                                   String,
  EAction                                     String,
  ELabel                                      String,
  EValue                                      UInt64,

  --- utm
  UTMValid                                    UInt8, -- bool
  UTMExist                                    UInt8, -- bool
  UTMSource                                   String,
	UTMMedium                                   String,
	UTMCampaign                                 String,
	UTMID                                       String,
	UTMTerm                                     String,
	UTMContent                                  String,

  -- user agent
	UaType                                      LowCardinality(String),
	UaFull                                      String,
	UaChecksum                                  FixedString(40),
	UaBrowserName                               LowCardinality(String),
	UaBrowserVersionMajor                       UInt64,
	UaBrowserVersion                            String,
	UaOSName                                    LowCardinality(String),
	UaOSVersionMajor                            UInt64,
	UaOSVersion                                 String,
	UaDeviceBrand                               LowCardinality(String),
	UaDeviceFamily                              String,
	UaDeviceModel                               String,

  -- screen
	ScrScreenOrientation                        UInt8, -- bool
	ScrScreenOrientationIsPortrait              UInt8, -- bool
	ScrScreenOrientationIsSecondary             UInt8, -- bool
	ScrScreen                                   String,
	ScrScreenWidth                              UInt16,
	ScrScreenHeight                             UInt16,
	ScrViewport                                 String,
	ScrViewportWidth                            UInt16,
	ScrViewportHeight                           UInt16,
	ScrResoluton                                String,
	ScrResolutonWidth                           UInt16,
	ScrResolutonHeight                          UInt16,
	ScrDevicePixelRatio                         Float64,
	ScrColorDepth                               UInt8,

  -- geo:asn
  GeoIPAutonomousSystemNumber                 UInt16,
  GeoIPAutonomousSystemOrganization           String,
  -- geo:ip
  GeoIPAdministratorArea                      String,
  GeoIPCity                                   String,
  GeoIPCityGeoNameID                          UInt16,
  GeoIPCountry                                LowCardinality(String),
  GeoIPLocationLatitude                       Float64,
  GeoIPLocationLongitude                      Float64,
  -- geo:ip+
  GeoIPLocation                               Point,
  -- geo:client
	GeoClientAdministratorArea                  String,
	GeoClientCity                               String,
	GeoClientCityGeoNameID                      UInt16,
	GeoClientCountry                            String,
	GeoClientLocationLatitude                   Float64,
	GeoClientLocationLongitude                  Float64,
  -- geo:client+
	GeoClientLocation                           Point,
  -- geo:result
	GeoResultAdministratorArea                  String,
	GeoResultCity                               String,
	GeoResultCityGeoNameID                      UInt16,
	GeoResultCountry                            String,
	GeoResultFromClient                         UInt8, -- bool
	GeoResultLocationLatitude                   Float64,
	GeoResultLocationLongitude                  Float64,
  -- geo:result+
	GeoResultLocation                           Point,

  -- client
  CIDType                                     UInt8,
  CIDUserChecksum                             FixedString(40),
  CIDSessionChecksum                          FixedString(40),
  CIDStdInitTime                              Datetime,
  CIDStdSessionTime                           Datetime,

  -- requirements
  IP                                          IPv4,
	PublicInstaceID                             String,
	Mode                                        UInt8,
  CursorID                                    UInt64,
	Created                                     Datetime
)
ENGINE = MergeTree()
ORDER BY (Created, xxHash32(CIDUserChecksum))
PARTITION BY toYYYYMM(Created)
SAMPLE BY (xxHash32(CIDUserChecksum));
