package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"net"
	"net/url"
	"time"
)

type breadCrumb struct {
	BCIsProcessed bool
	BCLevel       uint8
	BCN1          string
	BCN2          string
	BCN3          string
	BCN4          string
	BCN5          string
	BCP1          string
	BCP2          string
	BCP3          string
	BCP4          string
	BCP5          string
}

type performance struct {
	PerfIsProcessed        bool
	PerfPageLoadTime       uint16
	PerfDomainLookupTime   uint16
	PerfTCPConnectTime     uint16
	PerfServerResponseTime uint16
	PerfPageDownloadTime   uint16
	PerfRedirectTime       uint16
	PerfDOMInteractiveTime uint16
	PerfContentLoadTime    uint16
	PerfResource           uint16
}

type recordEvent struct {
	ECategory string
	EAction   string
	ELabel    string
	EValue    uint64
}

type record struct {
	ClientErrorMessage string
	ClientErrorObject  string
	Created            time.Time
	CursorID           uint64
	Mode               uint8
	modeString         string
	PublicInstaceID    string
	IP                 net.IP

	PEntityID         string
	PEntityModule     string
	PEntityTaxonomyID string
	PURL              *url.URL
	PCanonicalURL     *url.URL
	PTitle            string
	PLang             string
	PIsIframe         bool
	PIsTouchSupport   bool
	PKeywords         []string
	PRefererURL       refererData
	SRefererURL       refererData

	CID clientID

	EventCount int
	Events     []recordEvent

	UserAgentResult userAgentResult
	Utm             utm
	GeoResult       geoResult
	ScreenInfo      screenInfo
	BreadCrumb      breadCrumb
	Performance     performance

	UserIDOrName string
}

const (
	// PageView from 1-99
	recordModePageViewJavaScript    uint8 = 1
	recordModePageViewImageLegacy   uint8 = 2
	recordModePageViewImageNoScript uint8 = 3
	recordModePageViewAMP           uint8 = 4
	recordModePageViewAMPImage      uint8 = 5
	// rest of type must be added incremental

	// PageView from 100-199
	recordModeEventOther         uint8 = 100
	recordModeEventJSInPageView  uint8 = 101
	recordModeEventServiceWorker uint8 = 102
	recordModeEventAPI           uint8 = 103
	recordModeEventJSCross       uint8 = 104
	// rest of type must be added incremental

	// client error
	recordModeClientError uint8 = 255
)

var recordModeMap map[string]uint8

func init() {
	recordModeMap = make(map[string]uint8)

	// pageview
	recordModeMap["pv_js"] = recordModePageViewJavaScript
	recordModeMap["pv_il"] = recordModePageViewImageLegacy
	recordModeMap["pv_ins"] = recordModePageViewImageNoScript
	recordModeMap["pv_amp"] = recordModePageViewAMP
	recordModeMap["pv_amp_i"] = recordModePageViewAMPImage

	// event
	recordModeMap["e_o"] = recordModeEventOther
	recordModeMap["e_js_pv"] = recordModeEventJSInPageView
	recordModeMap["e_js_c"] = recordModeEventJSCross
	recordModeMap["e_sw"] = recordModeEventServiceWorker
	recordModeMap["e_api"] = recordModeEventAPI

	// etc
	recordModeMap["err"] = recordModeClientError
}

func validateMode(m string) (uint8, error) {
	if value, ok := recordModeMap[m]; ok {
		return value, nil
	}
	return 0, errors.New("invalid mode of request")
}

func newRecord(modeQuery string, publicInstaceIDQuery string) (*record, error) {
	r := record{
		Created: time.Now(),
	}

	mode, err := validateMode(modeQuery)
	if err != nil {
		return &r, err
	}

	publicInstaceID, err := validatePublicInstaceID(publicInstaceIDQuery)
	if err != nil {
		return &r, err
	}

	r.modeString = modeQuery
	r.Mode = mode
	r.PublicInstaceID = publicInstaceID

	return &r, nil
}

func (r *record) setAPI(
	postData *postRequest,
) *errorMessage {
	if postData.API == nil {
		return &errorAPIFieldsAreMissing
	}

	n := time.Now().Unix()
	if postData.API.ClientTime <= n && postData.API.ClientTime >= (n-28800) {
		r.Created = time.Unix(postData.API.ClientTime, 0)
	}

	ip := net.ParseIP(postData.API.ClientIP)
	if ip == nil {
		return &errorAPIClientIPNotValid
	}

	r.IP = ip

	if postData.API.ClientUserAgent == "" {
		return &errorAPIClientUserAgentNotValid
	}

	r.CID = clientIDFromOther([]string{r.IP.String(), postData.API.ClientUserAgent})

	return nil
}

func (r *record) isPageView() bool {
	if r.Mode == recordModePageViewJavaScript ||
		r.Mode == recordModePageViewAMP ||
		r.Mode == recordModePageViewImageLegacy ||
		r.Mode == recordModePageViewImageNoScript ||
		r.Mode == recordModePageViewAMPImage {
		return true
	}
	return false
}

func (r *record) isImage() bool {
	if r.Mode == recordModePageViewAMPImage ||
		r.Mode == recordModePageViewImageLegacy ||
		r.Mode == recordModePageViewImageNoScript {
		return true
	}
	return false
}

func (r *record) verify(
	projectsManager *projects,
	privateKey string,
) *errorMessage {
	// initialize
	if r.Mode < 1 || r.IP == nil {
		return &errorRecordNotProccessedYet
	}

	// page view require matched project id and url of page view
	if r.isPageView() {
		if r.PURL == nil {
			return &errorURLRequiredAndMustBeValid
		}
		if !projectsManager.validateIDAndURL(r.PublicInstaceID, r.PURL) {
			return &errorProjectPublicIDAndURLDidNotMatched
		}
		return nil
	}

	// in api mode private key must matched
	if r.Mode == recordModeEventAPI && !projectsManager.validateIDAndPrivate(r.PublicInstaceID, privateKey) {
		return &errorAPIPrivateKeyFailed
	}

	// in page js event must match with page url
	if r.Mode == recordModeEventJSInPageView && !projectsManager.validateIDAndURL(r.PublicInstaceID, r.PURL) {
		return &errorProjectPublicIDAndURLDidNotMatched
	}

	if r.Mode == recordModeClientError && r.PURL != nil && !projectsManager.validateIDAndURL(r.PublicInstaceID, r.PURL) {
		return &errorProjectPublicIDAndURLDidNotMatched
	}

	if r.Mode > 99 && r.EventCount < 1 {
		return &errorEventsAreEmpty
	}

	return nil
}

func (r *record) setReferer(refererParser *refererParser, refererURL *url.URL) {
	r.PRefererURL = refererParser.parse(r.PURL, refererURL)
}

func (r *record) setQueryParameters(
	qURL string,
	qCanonical string,
	qTitle string,
	qLang string,
	qEntityID string,
	qEntityModule string,
	qEntityTaxonomyID string,
) {
	r.PURL = getURL(qURL)
	r.PCanonicalURL = getURL(qCanonical)
	r.PTitle = qTitle
	r.PLang = sanitizeLanguage(qLang)
	r.PEntityID = sanitizeEntityID(qEntityID)
	r.PEntityModule = sanitizeName(qEntityModule)
	r.PEntityTaxonomyID = sanitizeEntityTaxonomyID(qEntityTaxonomyID)
}

func (r *record) setPostRequest(
	postRequest *postRequest,
	refererParser *refererParser,
	geoParser *geoParser,
) {
	if postRequest.ClientErrorMessage != "" {
		r.ClientErrorMessage = postRequest.ClientErrorMessage
		r.ClientErrorObject = postRequest.ClientErrorObject
	}

	if postRequest.Page != nil {
		r.PURL = getURL(postRequest.Page.URL)
		r.UserIDOrName = postRequest.Page.UserIDOrName
		r.PCanonicalURL = getURL(postRequest.Page.CanonicalURL)
		r.PTitle = postRequest.Page.Title
		r.PLang = sanitizeLanguage(postRequest.Page.Lang)
		r.PEntityID = sanitizeEntityID(postRequest.Page.MainEntityID)
		r.PEntityModule = sanitizeName(postRequest.Page.MainEntityModule)
		r.PEntityTaxonomyID = sanitizeEntityTaxonomyID(postRequest.Page.MainEntityTaxonomyID)

		if postRequest.Page.ScreenSize != "" &&
			postRequest.Page.ViewportSize != "" {
			r.ScreenInfo = parseScreenSize(
				postRequest.Page.ScreenOrientationType,
				postRequest.Page.ScreenSize,
				postRequest.Page.ViewportSize,
				postRequest.Page.DevicePixelRatio,
				postRequest.Page.ColorDepth,
			)
		}

		if postRequest.Page.RefererURL != "" {
			r.PRefererURL = refererParser.parse(r.PURL, getURL(postRequest.Page.RefererURL))
		}

		if postRequest.Page.RefererSessionURL != "" {
			r.SRefererURL = refererParser.parse(r.PURL, getURL(postRequest.Page.RefererSessionURL))
		}

		if postRequest.Page.PerformanceData != nil {
			r.Performance.PerfIsProcessed = true
			r.Performance.PerfPageLoadTime = uint16FromString(postRequest.Page.PerformanceData.PerfPageLoadTime)
			r.Performance.PerfDomainLookupTime = uint16FromString(postRequest.Page.PerformanceData.PerfDomainLookupTime)
			r.Performance.PerfTCPConnectTime = uint16FromString(postRequest.Page.PerformanceData.PerfTCPConnectTime)
			r.Performance.PerfServerResponseTime = uint16FromString(postRequest.Page.PerformanceData.PerfServerResponseTime)
			r.Performance.PerfPageDownloadTime = uint16FromString(postRequest.Page.PerformanceData.PerfPageDownloadTime)
			r.Performance.PerfRedirectTime = uint16FromString(postRequest.Page.PerformanceData.PerfRedirectTime)
			r.Performance.PerfDOMInteractiveTime = uint16FromString(postRequest.Page.PerformanceData.PerfDOMInteractiveTime)
			r.Performance.PerfContentLoadTime = uint16FromString(postRequest.Page.PerformanceData.PerfContentLoadTime)
			r.Performance.PerfResource = postRequest.Page.PerformanceData.PerfResource
		}

		if postRequest.Page.PageBreadcrumbObject != nil && r.PURL != nil {
			u1 := getURL(postRequest.Page.PageBreadcrumbObject.U1)
			if u1 != nil && u1.Host == r.PURL.Host && postRequest.Page.PageBreadcrumbObject.N1 != "" {
				r.BreadCrumb.BCIsProcessed = true
				r.BreadCrumb.BCLevel = 1
				r.BreadCrumb.BCN1 = postRequest.Page.PageBreadcrumbObject.N1
				r.BreadCrumb.BCP1 = getURLPath(getURL(postRequest.Page.PageBreadcrumbObject.U1))

				u2 := getURL(postRequest.Page.PageBreadcrumbObject.U2)
				if u2 != nil && u2.Host == r.PURL.Host && postRequest.Page.PageBreadcrumbObject.N2 != "" {
					r.BreadCrumb.BCLevel = 2
					r.BreadCrumb.BCN2 = postRequest.Page.PageBreadcrumbObject.N2
					r.BreadCrumb.BCP2 = getURLPath(getURL(postRequest.Page.PageBreadcrumbObject.U2))

					u3 := getURL(postRequest.Page.PageBreadcrumbObject.U3)
					if u3 != nil && u3.Host == r.PURL.Host && postRequest.Page.PageBreadcrumbObject.N3 != "" {
						r.BreadCrumb.BCLevel = 3
						r.BreadCrumb.BCN3 = postRequest.Page.PageBreadcrumbObject.N3
						r.BreadCrumb.BCP3 = getURLPath(getURL(postRequest.Page.PageBreadcrumbObject.U3))

						u4 := getURL(postRequest.Page.PageBreadcrumbObject.U4)
						if u4 != nil && u4.Host == r.PURL.Host && postRequest.Page.PageBreadcrumbObject.N4 != "" {
							r.BreadCrumb.BCLevel = 4
							r.BreadCrumb.BCN4 = postRequest.Page.PageBreadcrumbObject.N4
							r.BreadCrumb.BCP4 = getURLPath(getURL(postRequest.Page.PageBreadcrumbObject.U4))

							u5 := getURL(postRequest.Page.PageBreadcrumbObject.U5)
							if u5 != nil && u5.Host == r.PURL.Host && postRequest.Page.PageBreadcrumbObject.N5 != "" {
								r.BreadCrumb.BCLevel = 5
								r.BreadCrumb.BCN5 = postRequest.Page.PageBreadcrumbObject.N5
								r.BreadCrumb.BCP5 = getURLPath(getURL(postRequest.Page.PageBreadcrumbObject.U5))
							}
						}
					}

				}
			}
		}

		if postRequest.Page.GeographyData != nil {
			r.GeoResult = geoParser.clientLocationUpdate(
				r.GeoResult,
				postRequest.Page.GeographyData.Lat,
				postRequest.Page.GeographyData.Lon,
			)
		}

		r.PKeywords = parseKeywords(postRequest.Page.PageKeywords)
		r.PIsIframe = postRequest.Page.IsIframe
		r.PIsTouchSupport = postRequest.Page.IsTouchSupport
	}

	if postRequest.Events != nil {
		events := []recordEvent{}
		for _, ev := range *postRequest.Events {
			category := sanitizeName(ev.Category)
			action := sanitizeName(ev.Action)

			if category != "" && action != "" {
				re := recordEvent{
					ECategory: category,
					EAction:   action,
					ELabel:    ev.Label,
					EValue:    ev.Value,
				}
				events = append(events, re)
			}

		}
		r.Events = events
		r.EventCount = len(events)
	}

	if r.Mode == recordModePageViewJavaScript || r.Mode == recordModeEventJSInPageView && postRequest.CIDStd != "" {
		cid, cidErr := clientIDStandardParser(postRequest.CIDStd)
		if cidErr == nil {
			r.CID = cid
		}
	} else if r.Mode == recordModePageViewAMP && postRequest.CIDAmp != "" {
		r.CID = clientIDFromAMP(postRequest.CIDAmp)
	} else if r.Mode == recordModeClientError && r.IP != nil {
		r.CID = clientIDFromOther([]string{r.IP.String(), r.UserAgentResult.UaFull})
	}
}

func (r *record) finalize() ([]byte, error) {
	if r.Mode < 1 || !r.CID.Valid {
		return nil, errors.New("mode not processed or missing cid")
	}

	if r.PURL != nil {
		r.Utm = parseUTM(r.PURL)
	}

	if r.isPageView() {
		r.CursorID = getCursorID()
	}

	buf := &bytes.Buffer{}
	if err := gob.NewEncoder(buf).Encode(*r); err != nil {
		return nil, err
	}

	defer promMetricRecordMode.WithLabelValues(r.modeString).Inc()

	return buf.Bytes(), nil
}
