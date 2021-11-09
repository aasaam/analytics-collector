package main

import (
	"net"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	// Page view events is under 50
	RecordTypePageViewJS            uint8 = 0
	RecordTypePageViewLegacyImage   uint8 = 1
	RecordTypePageViewNoScriptImage uint8 = 2
	RecordTypePageViewAMP           uint8 = 3
	RecordTypePageViewAMPNoScript   uint8 = 4
	RecordTypePageViewMobile        uint8 = 5

	// Event
	RecordTypeEventJS            uint8 = 50
	RecordTypeEventServiceWorker uint8 = 51
	RecordTypeEventAPI           uint8 = 52
)

// Record is finalize record ready for insert in database
type Record struct {
	RecordType uint8
	Time       time.Time
	ProjectID  string
	Cursor     uint64

	GeoIP            GeoResult
	UserAgent        UserAgentResult
	UTM              UTM
	ClientIdentifier ClientIdentifier
	ScreenInfo       ScreenInfo
	RefererData      RefererData
}

// RecordTypeMap is map from URL get params to unit8 value
var RecordTypeMap map[string]uint8

// PostEventItem is single event data
type PostEventItem struct {
	Category string `json:"ec"`
	Action   string `json:"ea"`
	Label    string `json:"el,omitempty"`
	Value    int    `json:"ev,omitempty"`
}

// PostRecordsPage is page view data
type PostRecordsPage struct {
	Screen     string   `json:"scr,omitempty"`
	Viewport   string   `json:"vps,omitempty"`
	ColorDepth string   `json:"cd,omitempty"`
	Title      string   `json:"t,omitempty"`
	Keywords   []string `json:"k,omitempty"`
	User       string   `json:"usr,omitempty"`

	DevicePixelRatio  string `json:"dpr,omitempty"`
	IFrame            bool   `json:"if"`
	ScreenOrientation string `json:"so,omitempty"`
	RefererSate       string `json:"rs,omitempty"`

	Performance PostRecordsPagePerformance `json:"prf,omitempty"`
	ClientGeo   PostRecordsPageGeo         `json:"geo,omitempty"`
}

// PostRecordsPagePerformance is performance data
type PostRecordsPagePerformance struct {
	PageLoadTime       string `json:"plt,omitempty"`
	DomainLookupTime   string `json:"dlt,omitempty"`
	TcpConnectTime     string `json:"tct,omitempty"`
	ServerResponseTime string `json:"srt,omitempty"`
	PageDownloadTime   string `json:"pdt,omitempty"`
	RedirectTime       string `json:"rt,omitempty"`
	DomInteractiveTime string `json:"dit,omitempty"`
	ContentLoadTime    string `json:"clt,omitempty"`
}

// PostRecordsPageGeo is geo location of client
type PostRecordsPageGeo struct {
	ClientLatitude       float64 `json:"lat"`
	ClientLongitude      float64 `json:"lon"`
	ClientAccuracyRadius float64 `json:"acc"`
}

// PostRecordsData is incoming request for submit new records
type PostRecordsData struct {
	ClientIdentifier string          `json:"c,omitempty"`
	PageView         PostRecordsPage `json:"p,omitempty"`
	Events           []PostEventItem `json:"e,omitempty"`
}

// GetRecordsData
type GetRecordsData struct {
	Mode               uint8
	ProjectPublicHash  string
	ProjectPrivateHash string
	ClientIP           *net.IP
	ClientUserAgent    string
	URL                *url.URL
	Referer            *url.URL
	Canonical          *url.URL
	MainID             string
}

func init() {
	RecordTypeMap = make(map[string]uint8)
	RecordTypeMap["jsp"] = RecordTypePageViewJS
	RecordTypeMap["il"] = RecordTypePageViewLegacyImage
	RecordTypeMap["ins"] = RecordTypePageViewNoScriptImage
	RecordTypeMap["amp"] = RecordTypePageViewAMP
	RecordTypeMap["ia"] = RecordTypePageViewAMPNoScript
	RecordTypeMap["jse"] = RecordTypeEventJS
	RecordTypeMap["sw"] = RecordTypeEventServiceWorker
	RecordTypeMap["api"] = RecordTypeEventAPI
}

func NewRecordsFromRequest(
	projects *Projects,
	referrerParser *ReferrerParser,
	geoParser *GeoParser,
	c *fiber.Ctx,
) ([]Record, error) {
	var result []Record
	return result, nil
}

// // RecordData is whole data from client
// type RecordData struct {
// 	APIClientIP         string            `json:"c_ip,omitempty"`
// 	APIClientUserAgent  string            `json:"c_ua,omitempty"`
// 	StdClientIdentifier string            `json:"c,omitempty"`
// 	Page                RecordDataPage    `json:"p,omitempty"`
// 	Events              []RecordDataEvent `json:"e,omitempty"`
// }

// type Record struct {
// 	// one of: amp,img,pageview,event
// 	Type      uint8
// 	Time      time.Time
// 	ProjectID string
// 	Cursor    uint64

// 	GeoIP            GeoIPResult
// 	UserAgent        UserAgentResult
// 	UTM              UTM
// 	ClientIdentifier ClientIdentifier
// 	ScreenInfo       ScreenInfo
// 	RefererData      RefererData

// 	Title     string
// 	EntityID  string
// 	User      string
// 	Keywords  []string
// 	Canonical string
// }

// var hexDecRegexReplace = regexp.MustCompile(`[^0-9]`)

// func isEventRecord(recordType uint8) bool {
// 	if recordType >= 10 && recordType < 50 {
// 		return true
// 	}
// 	return false
// }

// func isPageViewRecord(recordType uint8) bool {
// 	if recordType < 10 {
// 		return true
// 	}
// 	return false
// }

// func cursorProcess(projectId string) uint64 {
// 	micro := strconv.FormatInt(time.Now().UnixMicro(), 10)

// 	h := sha1.New()
// 	h.Write([]byte(projectId))
// 	hex := hex.EncodeToString(h.Sum(nil))
// 	hexDec := hexDecRegexReplace.ReplaceAllString(hex, "")

// 	n, e := strconv.ParseInt(micro+hexDec[0:3], 10, 64)
// 	if e != nil {
// 		panic(e)
// 	}
// 	return uint64(n)
// }

// func NewRecord(projectId string, recordType uint8) *Record {
// 	r := Record{}
// 	r.Type = recordType
// 	r.ProjectID = projectId
// 	r.Time = time.Now()
// 	if isPageViewRecord(recordType) {
// 		r.Cursor = cursorProcess(projectId)
// 	}

// 	return &r
// }

// func NewRecordsFromRequest(
// 	projects *Projects,
// 	c *fiber.Ctx,
// ) error {
// 	mode, modeExist := RecordTypeMap[c.Get("m", "")]
// 	if !modeExist {
// 		return c.Status(fiber.StatusBadRequest).SendString("invalid request mode")
// 	}

// 	url, err := url.Parse(c.Get("u", ""))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("invalid url")
// 		return err
// 	}

// 	referer, err := url.Parse(c.Get("u", ""))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("invalid url")
// 		return err
// 	}

// 	if mode == RecordTypeEventAPI {

// 	}

// 	// service worker

// 	fmt.Println(mode)
// 	return nil
// 	// validationResult =
// 	// if mode == RecordTypeEventAPI {

// 	// } else {

// 	// }

// 	// if (mode )
// 	// projectPublicHash := c.Get("i", "")
// 	// projects.Validate()
// 	// return nil
// }
