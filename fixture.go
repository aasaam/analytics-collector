package main

import (
	"encoding/base64"
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
	"gopkg.in/yaml.v2"
)

type fixture struct {
	ECategory       []string            `yaml:"ECategory"`
	EAction         []string            `yaml:"EAction"`
	EIdent          []string            `yaml:"EIdent"`
	Created         []bool              `yaml:"Created"`
	PIsIframe       []bool              `yaml:"PIsIframe"`
	PIsTouchSupport []bool              `yaml:"PIsTouchSupport"`
	PLang           []string            `yaml:"PLang"`
	PEntityModule   []string            `yaml:"PEntityModule"`
	Referer         []string            `yaml:"Referer"`
	Titles          map[string][]string `yaml:"Titles"`
	Geo             []struct {
		Country string     `yaml:"Country"`
		Lat     [2]float64 `yaml:"Lat"`
		Lon     [2]float64 `yaml:"Lon"`
	} `yaml:"Geo"`
	Mode              []string            `yaml:"Mode"`
	UserAgent         []string            `yaml:"UserAgent"`
	PEntityTaxonomyID []uint16            `yaml:"PEntityTaxonomyID"`
	PublicInstanceID  map[string][]string `yaml:"PublicInstanceID"`
}

func (f *fixture) stdCID() string {
	initTime := time.Now().Add(time.Duration(fake.Number(-80, -40)) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(fake.Number(-20, -10)) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	return base64.StdEncoding.EncodeToString([]byte(cid))
}

func (f *fixture) geoResult(ip string) geoResult {
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

func (f *fixture) rand(min int, max int) int {
	return fake.Number(min, max)
}

func (f *fixture) record(
	refererParser *refererParser,
	userAgentParser *userAgentParser,
) *record {
	past, _ := time.Parse("2006-01-02T15:04:05.000Z", "2020-01-01T00:00:00.000Z")

	PublicInstanceID := reflect.ValueOf(f.PublicInstanceID).MapKeys()

	r := record{}

	r.Created = time.Now()
	usePast := false
	if f.Created[rand.Intn(len(f.Created))] {
		r.Created = time.Now()
	} else {
		usePast = true
		r.Created = fake.DateRange(past, time.Now().Add(-1*time.Minute))
	}

	ipString := fake.IPv4Address()
	r.CID = clientIDNoneSTD([]string{ipString}, clientIDTypeOther)
	r.IP = net.ParseIP(ipString)

	r.GeoResult = f.geoResult(ipString)
	r.UserAgentResult = userAgentParser.parse(f.UserAgent[rand.Intn(len(f.UserAgent))])
	r.PublicInstanceID = PublicInstanceID[rand.Intn(len(PublicInstanceID))].String()

	var modeErr error
	r.Mode, modeErr = validateMode(f.Mode[rand.Intn(len(f.Mode))])
	if modeErr != nil {
		panic(modeErr)
	}

	if r.Mode < 100 || r.Mode == recordModeEventJSInPageView {
		if r.Mode == recordModePageViewJavaScript {
			cid, cidErr := clientIDStandardParser(f.stdCID())
			if cidErr == nil {
				r.CID = cid
			} else {
				panic(cidErr)
			}
		}

		lang := f.PLang[rand.Intn(len(f.PLang))]
		title := fake.Sentence(fake.Number(5, 30))
		if lang != "en" {
			Titles, foundTitles := f.Titles[lang]
			if foundTitles {
				title = Titles[rand.Intn(len(Titles))]
			}
		}
		domains := f.PublicInstanceID[r.PublicInstanceID]
		domain := domains[rand.Intn(len(domains))]

		r.PIsIframe = f.PIsIframe[rand.Intn(len(f.PIsIframe))]
		r.PIsTouchSupport = f.PIsTouchSupport[rand.Intn(len(f.PIsTouchSupport))]
		r.PKeywords = strings.Split(title, " ")

		u, _ := url.Parse(fake.URL())
		u.Host = domain
		u.Scheme = "https"

		cu := u.String()

		if fake.Number(0, 4) == 4 {
			values := u.Query()
			values.Set("utm_source", fake.Noun())
			values.Set("utm_medium", fake.Noun())
			values.Set("utm_campaign", fake.Noun())
			if fake.Number(0, 2) == 2 {
				values.Set("utm_id", fake.Noun())
				values.Set("utm_term", fake.Noun())
				values.Set("utm_content", fake.Noun())
			}

			u.RawQuery = values.Encode()
		}
		entID := ""
		entMod := ""
		var entTID uint16 = 0
		if fake.Number(0, 5) == 5 {
			entID = strconv.Itoa(fake.Number(1, 100000))
			entMod = f.PEntityModule[rand.Intn(len(f.PEntityModule))]
			entTID = f.PEntityTaxonomyID[rand.Intn(len(f.PEntityTaxonomyID))]
		}

		r.setQueryParameters(
			u.String(),
			cu,
			title,
			lang,
			entID,
			entMod,
			strconv.Itoa(int(entTID)),
		)

		// Referer
		if r.pURL != nil {

			// session
			fakeReferer1 := f.Referer[rand.Intn(len(f.Referer))]

			if fakeReferer1 != "" {
				ref, refErr := url.Parse(fakeReferer1)
				if refErr == nil {
					r.SRefererURL = refererParser.parse(r.pURL, ref)
				}
			}

			// referer
			fakeReferer2 := f.Referer[rand.Intn(len(f.Referer))]
			if fakeReferer2 != "" {
				ref, refErr := url.Parse(fakeReferer2)
				if refErr == nil {
					r.PRefererURL = refererParser.parse(r.pURL, ref)
				}
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
				ev.EIdent = f.EIdent[rand.Intn(len(f.EIdent))]
			}
			if fake.Number(0, 3) == 3 {
				ev.EValue = uint64(fake.Number(1, 10000))
			}

			Events = append(Events, ev)
		}
		r.Events = Events
	}

	if r.pURL != nil {
		r.Utm = parseUTM(r.pURL)
	}

	// finalize
	if r.isPageView() && !usePast {
		r.CursorID, _ = getCursorID()
	}

	return &r
}

func fixtureLoad(path string) (*fixture, error) {
	yamlData, yamlDataErr := os.ReadFile(path)
	if yamlDataErr != nil {
		return nil, yamlDataErr
	}
	f := fixture{}
	yamlParseErr := yaml.Unmarshal(yamlData, &f)
	if yamlParseErr != nil {
		return nil, yamlParseErr
	}
	return &f, nil
}
