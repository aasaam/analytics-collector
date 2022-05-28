INSERT INTO ClientErrors (
  Msg                                         , -- String,
  Err                                         , -- String,

  PURL                                        , -- String,
  PURLChecksum                                , -- String,

  -- geo:asn
  GeoIPAutonomousSystemNumber                 , -- UInt16,
  GeoIPAutonomousSystemOrganization           , -- String,
  -- geo:result
	GeoResultAdministratorArea                  , -- String,
	GeoResultCity                               , -- String,
	GeoResultCityGeoNameID                      , -- UInt32,
	GeoResultCountry                            , -- String,
	GeoResultLocationLatitude                   , -- Float64,
	GeoResultLocationLongitude                  , -- Float64,
  -- geo:result+
	GeoResultLocation                           , -- Point,

  -- user agent
	UaType                                      , -- LowCardinality(String),
	UaFull                                      , -- String,
	UaChecksum                                  , -- FixedString(40),
	UaBrowserName                               , -- LowCardinality(String),
	UaBrowserVersionMajor                       , -- UInt64,
	UaBrowserVersion                            , -- String,
	UaOSName                                    , -- LowCardinality(String),
	UaOSVersionMajor                            , -- UInt64,
	UaOSVersion                                 , -- String,
	UaDeviceBrand                               , -- LowCardinality(String),
	UaDeviceFamily                              , -- String,
	UaDeviceModel                               , -- String,

  -- requirements
  IP                                          , -- IPv4,
	PublicInstanceID                             , -- String,
	Mode                                        , -- UInt8,
	Created                                     -- Datetime
)
