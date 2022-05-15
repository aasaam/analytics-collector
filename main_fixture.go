package main

import (
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type fixture struct {
	ECategory       []string `yaml:"ECategory"`
	EAction         []string `yaml:"EAction"`
	Created         []bool   `yaml:"Created"`
	PIsIframe       []bool   `yaml:"PIsIframe"`
	PIsTouchSupport []bool   `yaml:"PIsTouchSupport"`
	PEntityModule   []string `yaml:"PEntityModule"`
	Geo             []struct {
		Country string     `yaml:"Country"`
		Lat     [2]float64 `yaml:"Lat"`
		Lon     [2]float64 `yaml:"Lon"`
	} `yaml:"Geo"`
	Mode              []string            `yaml:"Mode"`
	UserAgent         []string            `yaml:"UserAgent"`
	PEntityTaxonomyID []string            `yaml:"PEntityTaxonomyID"`
	PublicInstanceID  map[string][]string `yaml:"PublicInstanceID"`
}

func fakeStdCID() string {
	initTime := time.Now().Add(time.Duration(fake.Number(-40, -80)) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(fake.Number(-10, -20)) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	return base64.StdEncoding.EncodeToString([]byte(cid))
}

func fakeGeoResult(ip string, f *fixture) geoResult {
	r := geoResult{
		GeoIsProcessed: false,
	}

	r.GeoIP = ip
	r.GeoIPAutonomousSystemNumber = uint16(fake.Number(100, 60000))
	r.GeoIPAutonomousSystemOrganization = "Fake ISP " + fake.Noun()

	geo := f.Geo[rand.Intn(len(f.Geo))]
	r.GeoIPAdministratorArea = fake.City()
	r.GeoIPCity = fake.City()
	r.GeoIPCityGeoNameID = uint32(fake.Number(100, 60000))
	r.GeoIPCountry = geo.Country
	var round float64 = 100
	r.GeoIPLocationLatitude = math.Round(fake.Float64Range(geo.Lat[0], geo.Lat[1])*round) / round
	r.GeoIPLocationLongitude = math.Round(fake.Float64Range(geo.Lon[0], geo.Lon[1])*round) / round

	r.GeoResultAdministratorArea = r.GeoIPAdministratorArea
	r.GeoResultCity = r.GeoIPCity
	r.GeoResultCityGeoNameID = r.GeoIPCityGeoNameID
	r.GeoResultCountry = r.GeoIPCountry
	r.GeoResultFromClient = false
	r.GeoResultLocationLatitude = r.GeoIPLocationLatitude
	r.GeoResultLocationLongitude = r.GeoIPLocationLongitude

	if fake.Number(0, 5) == 5 {
		var round2 float64 = 100000
		r.GeoIPLocationLatitude = math.Round(fake.Float64Range(geo.Lat[0], geo.Lat[1])*round2) / round2
		r.GeoIPLocationLongitude = math.Round(fake.Float64Range(geo.Lon[0], geo.Lon[1])*round2) / round2
		r.GeoResultFromClient = true
		r.GeoClientAdministratorArea = r.GeoIPAdministratorArea
		r.GeoClientCity = r.GeoIPCity
		r.GeoClientCityGeoNameID = r.GeoIPCityGeoNameID
		r.GeoClientCountry = r.GeoIPCountry
		r.GeoClientLocationLatitude = r.GeoIPLocationLatitude
		r.GeoClientLocationLongitude = r.GeoIPLocationLongitude
	}

	return r
}

func runFixture(c *cli.Context) error {

	conf := newConfig(
		"debug",
		0,
		true,
		"",
		"",
	)

	fixtureInterval := time.Duration(c.Int("fixture-interval")) * time.Second

	yamlData, yamlDataErr := os.ReadFile(c.String("fixture-yaml-path"))
	if yamlDataErr != nil {
		panic("cannot load yaml file: " + yamlDataErr.Error())
	}

clickHouseInitStep:

	_, _, clickhouseInitErr := getClickhouseConnection(
		c.String("clickhouse-servers"),
		c.String("clickhouse-database"),
		c.String("clickhouse-username"),
		c.String("clickhouse-password"),
		c.Int("clickhouse-max-execution-time"),
		c.Int("clickhouse-dial-timeout"),
		c.Bool("test-mode"),
		c.Bool("clickhouse-compression-lz4"),
		c.Int("clickhouse-max-idle-conns"),
		c.Int("clickhouse-max-open-conns"),
		c.Int("clickhouse-conn-max-lifetime"),
		c.Int("clickhouse-max-block-size"),
		nil,
		nil,
	)

	if clickhouseInitErr != nil {
		conf.getLogger().
			Error().
			Msg(clickhouseInitErr.Error())
		time.Sleep(fixtureInterval)
		goto clickHouseInitStep
	}

	f := fixture{}
	yamlParseErr := yaml.Unmarshal(yamlData, &f)
	if yamlParseErr != nil {
		panic("cannot parse yaml file:" + yamlParseErr.Error())
	}

	userAgentParser := newUserAgentParser()

	refererParser := newRefererParser()
	refStdDomains := reflect.ValueOf(refererParser.domainMap).MapKeys()

	for {

		rand.Seed(time.Now().Unix())

		records := make([]record, 0)
		past, _ := time.Parse("2006-01-02T15:04:05.000Z", "2020-01-01T00:00:00.000Z")

		PublicInstanceID := reflect.ValueOf(f.PublicInstanceID).MapKeys()

		for i := 0; i <= 1000; i++ {
			var r record
			var e error

			// created
			r.Created = time.Now()
			usePast := false
			if f.Created[rand.Intn(len(f.Created))] {
				r.Created = time.Now()
			} else {
				usePast = true
				r.Created = fake.DateRange(past, time.Now().Add(-1*time.Minute))
			}

			// geo
			ipString := fake.IPv4Address()
			r.CID = clientIDFromAMP(ipString)
			r.IP = net.ParseIP(ipString)
			if fake.Number(0, 5) > 3 {
				ref, _ := url.Parse(fake.URL())
				r.setReferer(refererParser, ref)
			}

			r.GeoResult = fakeGeoResult(ipString, &f)
			r.UserAgentResult = userAgentParser.parse(f.UserAgent[rand.Intn(len(f.UserAgent))])
			r.PublicInstanceID = PublicInstanceID[rand.Intn(len(PublicInstanceID))].String()

			r.Mode, e = validateMode(f.Mode[rand.Intn(len(f.Mode))])
			if e != nil {
				conf.getLogger().
					Error().
					Msg(e.Error())
				continue
			}

			// page vide
			if r.Mode < 100 || r.Mode == recordModeEventJSInPageView {
				if r.Mode == recordModePageViewJavaScript {
					cid, cidErr := clientIDStandardParser(fakeStdCID())
					if cidErr == nil {
						r.CID = cid
					} else {
						fmt.Println(cidErr)
					}
				}

				domains := f.PublicInstanceID[r.PublicInstanceID]
				domain := domains[rand.Intn(len(domains))]

				r.PIsIframe = f.PIsIframe[rand.Intn(len(f.PIsIframe))]
				r.PIsTouchSupport = f.PIsTouchSupport[rand.Intn(len(f.PIsTouchSupport))]
				r.PKeywords = strings.Split(fake.Sentence(5), " ")

				if fake.Number(0, 5) == 5 {
					r.UserIDOrName = fake.Username()
				}

				u, _ := url.Parse(fake.URL())
				u.Host = domain
				u.Scheme = "https"

				cu := u.String()

				if fake.Number(0, 5) == 5 {
					values := u.Query()
					values.Set("utm_source", fake.Noun())
					values.Set("utm_medium", fake.Noun())
					values.Set("utm_campaign", fake.Noun())
					if fake.Number(0, 5) == 5 {
						values.Set("utm_id", fake.Noun())
						values.Set("utm_term", fake.Noun())
						values.Set("utm_content", fake.Noun())
					}

					u.RawQuery = values.Encode()
				}

				entID := ""
				entMod := ""
				entTID := ""
				if fake.Number(0, 5) == 5 {
					entID = strconv.Itoa(fake.Number(1, 100000))
					entMod = f.PEntityModule[rand.Intn(len(f.PEntityModule))]
					entTID = f.PEntityTaxonomyID[rand.Intn(len(f.PEntityTaxonomyID))]
				}

				r.setQueryParameters(
					u.String(),
					cu,
					fake.Sentence(fake.Number(5, 30)),
					"en",
					entID,
					entMod,
					entTID,
				)

				// refere
				if r.PURL != nil && fake.Number(0, 3) == 3 {
					sRefURL, _ := url.Parse(fake.URL())
					if fake.Number(0, 5) == 5 {
						stdrefDomain := refStdDomains[rand.Intn(len(refStdDomains))].String()
						sRefURL.Host = stdrefDomain
					}
					r.SRefererURL = refererParser.parse(r.PURL, sRefURL)

					if fake.Number(0, 3) == 3 {
						refURL, _ := url.Parse(fake.URL())
						if fake.Number(0, 5) == 5 {
							stdrefDomain := refStdDomains[rand.Intn(len(refStdDomains))].String()
							refURL.Host = stdrefDomain
						}
						r.PRefererURL = refererParser.parse(r.PURL, refURL)
					}
				}
			}

			if r.Mode >= 100 && r.Mode < 200 { // it's event
				r.EventCount = fake.Number(0, 3)
				Events := make([]recordEvent, 0)
				for i := 0; i < r.EventCount; i += 1 {
					ev := recordEvent{
						ECategory: f.ECategory[rand.Intn(len(f.ECategory))],
						EAction:   f.EAction[rand.Intn(len(f.EAction))],
					}
					if fake.Number(0, 3) == 3 {
						ev.ELabel = fake.Noun()
					}
					if fake.Number(0, 3) == 3 {
						ev.EValue = uint64(fake.Number(1, 10000))
					}

					Events = append(Events, ev)
				}
				r.Events = Events
			}

			// finalize
			if r.PURL != nil {
				r.Utm = parseUTM(r.PURL)
			}

			// finalize
			if r.isPageView() && !usePast {
				r.CursorID = getCursorID()
			}

			records = append(records, r)
		}

		clickhouseConn, clickhouseCtx, clickhouseConnErr := getClickhouseConnection(
			c.String("clickhouse-servers"),
			c.String("clickhouse-database"),
			c.String("clickhouse-username"),
			c.String("clickhouse-password"),
			c.Int("clickhouse-max-execution-time"),
			c.Int("clickhouse-dial-timeout"),
			c.Bool("test-mode"),
			c.Bool("clickhouse-compression-lz4"),
			c.Int("clickhouse-max-idle-conns"),
			c.Int("clickhouse-max-open-conns"),
			c.Int("clickhouse-conn-max-lifetime"),
			c.Int("clickhouse-max-block-size"),
			nil,
			nil,
		)

		if clickhouseConnErr != nil {
			conf.getLogger().
				Error().
				Str("type", errorTypeApp).
				Str("on", "clickhouse-connection").
				Str("error", clickhouseConnErr.Error()).
				Send()
			time.Sleep(fixtureInterval)
			continue
		}

		recordsBatch, recordsBatchErr := clickhouseConn.PrepareBatch(
			clickhouseCtx, clickhouseInsertRecords,
		)

		if recordsBatchErr != nil {
			conf.getLogger().
				Error().
				Str("type", errorTypeApp).
				Str("on", "clickhouse-connection").
				Str("error", recordsBatchErr.Error()).
				Send()
			time.Sleep(fixtureInterval)
			continue
		}

		inserts := 0
		for _, rec := range records {

			if rec.EventCount > 0 {
				for i := 0; i < rec.EventCount; i++ {
					ECategory := rec.Events[i].ECategory
					EAction := rec.Events[i].EAction
					ELabel := rec.Events[i].ELabel
					EValue := rec.Events[i].EValue
					insertErr := insertRecordBatch(recordsBatch, rec, ECategory, EAction, ELabel, EValue)
					if insertErr != nil {
						conf.getLogger().
							Error().
							Str("type", errorTypeApp).
							Str("on", "record-insert").
							Str("error", insertErr.Error()).
							Send()
					}
					inserts += 1
				}
			} else {
				insertErr := insertRecordBatch(recordsBatch, rec, "", "", "", 0)
				if insertErr != nil {
					conf.getLogger().
						Error().
						Str("type", errorTypeApp).
						Str("on", "record-insert").
						Str("error", insertErr.Error()).
						Send()
				}
				inserts += 1
			}
		}

		recordsBatchSendErr := recordsBatch.Send()
		if recordsBatchSendErr != nil {
			conf.getLogger().
				Error().
				Str("type", errorTypeApp).
				Str("on", "record-batch-send").
				Str("error", recordsBatchSendErr.Error()).
				Send()
		}

		conf.getLogger().
			Warn().
			Int("inserted", inserts).
			Send()
		time.Sleep(fixtureInterval)
	}
}
