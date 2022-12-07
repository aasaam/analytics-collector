-- variables
SET allow_experimental_geo_types = 1;

-- database
CREATE DATABASE IF NOT EXISTS analytics;

-- ClientErrors
CREATE TABLE IF NOT EXISTS analytics.ClientErrors
(
  Msg                                         String,
  Err                                         String,

  PURL                                        String,
  PURLChecksum                                FixedString(24),

  -- geo:asn
  GeoIPAutonomousSystemNumber                 UInt16,
  GeoIPAutonomousSystemOrganization           String,
  -- geo:result
  GeoResultAdministratorArea                  String,
  GeoResultCity                               String,
  GeoResultCityGeoNameID                      UInt32,
  GeoResultCountry                            LowCardinality(String),
  GeoResultLocationLatitude                   Float64,
  GeoResultLocationLongitude                  Float64,
  -- geo:result+
  GeoResultLocation                           Point,

  -- user agent
  UaType                                      UInt8,
  UaFull                                      String,
  UaChecksum                                  FixedString(24),
  UaBrowserName                               LowCardinality(String),
  UaBrowserVersionMajor                       UInt64,
  UaBrowserVersion                            String,
  UaOSName                                    LowCardinality(String),
  UaOSVersionMajor                            UInt64,
  UaOSVersion                                 String,
  UaDeviceBrand                               LowCardinality(String),
  UaDeviceFamily                              String,
  UaDeviceModel                               String,

  -- requirements
  IP                                          IPv4,
  PublicInstanceID                            String,
  Mode                                        UInt8,
  Created                                     Datetime
)
ENGINE = ReplacingMergeTree()
ORDER BY (Created, xxHash32(PublicInstanceID))
PARTITION BY toYYYYMM(Created)
TTL Created + INTERVAL 1 WEEK;

-- Records
CREATE TABLE IF NOT EXISTS analytics.Records
(
  -- event
  ECategory                                   String,
  EAction                                     String,
  ELabel                                      String,
  EIdent                                      String,
  EValue                                      UInt64,

  -- custom segments
  Seg1Name                                    String,
  Seg2Name                                    String,
  Seg3Name                                    String,
  Seg4Name                                    String,
  Seg5Name                                    String,
  Seg1Value                                   String,
  Seg2Value                                   String,
  Seg3Value                                   String,
  Seg4Value                                   String,
  Seg5Value                                   String,

  -- page
  PIsIframe                                   UInt8, --bool
  PIsTouchSupport                             UInt8, --bool
  PURL                                        String,
  PURLChecksum                                FixedString(24),
  PTitle                                      String,
  PCanonicalURL                               String,
  PCanonicalURLChecksum                       FixedString(24),
  PLang                                       LowCardinality(String),
  PEntityID                                   String,
  PEntityModule                               String,
  PEntityTaxonomyID                           UInt16,
  PKeywords                                   Array(String),

  -- referer
  PRefererURLURL                              String,
  PRefererURLExternalHost                     String,
  PRefererURLExternalDomain                   String,
  PRefererURLExternalName                     String,
  PRefererURLScheme                           LowCardinality(String),
  PRefererURLExternalType                     UInt8,

  -- session referer
  SRefererURLURL                              String,
  SRefererURLExternalHost                     String,
  SRefererURLExternalDomain                   String,
  SRefererURLExternalName                     String,
  SRefererURLScheme                           LowCardinality(String),
  SRefererURLExternalType                     UInt8,

  --- utm
  UTMValid                                    UInt8, -- bool
  UTMExist                                    UInt8, -- bool
  UTMSource                                   String,
  UTMMedium                                   String,
  UTMCampaign                                 String,
  UTMID                                       String,
  UTMTerm                                     String,
  UTMContent                                  String,

  -- performance
  PerfIsProcessed                             UInt8, --bool
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
  BCLevel                                     UInt8,
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

  -- user agent
  UaType                                      UInt8,
  UaFull                                      String,
  UaChecksum                                  FixedString(24),
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
  ScrResolution                               String,
  ScrResolutionWidth                          UInt16,
  ScrResolutionHeight                         UInt16,
  ScrDevicePixelRatio                         Float64,
  ScrColorDepth                               UInt8,

  -- geo:asn
  GeoIPAutonomousSystemNumber                 UInt16,
  GeoIPAutonomousSystemOrganization           String,
  -- geo:ip
  GeoIPAdministratorArea                      String,
  GeoIPCity                                   String,
  GeoIPCityGeoNameID                          UInt32,
  GeoIPCountry                                LowCardinality(String),
  GeoIPLocationLatitude                       Float64,
  GeoIPLocationLongitude                      Float64,
  -- geo:ip+
  GeoIPLocation                               Point,
  -- geo:client
  GeoClientAdministratorArea                  String,
  GeoClientCity                               String,
  GeoClientCityGeoNameID                      UInt32,
  GeoClientCountry                            LowCardinality(String),
  GeoClientLocationLatitude                   Float64,
  GeoClientLocationLongitude                  Float64,
  -- geo:client+
  GeoClientLocation                           Point,
  -- geo:result
  GeoResultFromClient                         UInt8, -- bool
  GeoResultAdministratorArea                  String,
  GeoResultCity                               String,
  GeoResultCityGeoNameID                      UInt32,
  GeoResultCountry                            LowCardinality(String),
  GeoResultLocationLatitude                   Float64,
  GeoResultLocationLongitude                  Float64,
  -- geo:result+
  GeoResultLocation                           Point,

  -- client
  CidType                                     UInt8,
  CidUserChecksum                             FixedString(24),
  CidSessionChecksum                          FixedString(24),
  CidStdInitTime                              Datetime,
  CidStdSessionTime                           Datetime,

  -- requirements
  IP                                          IPv4,
  PublicInstanceID                            String,
  Mode                                        UInt8,
  CursorID                                    UInt64,
  Created                                     Datetime
)
ENGINE = ReplacingMergeTree()
ORDER BY (Created, Mode, xxHash32(PublicInstanceID))
PARTITION BY toYYYYMM(Created)
SAMPLE BY (xxHash32(PublicInstanceID));
