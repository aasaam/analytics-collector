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
	Created         time.Time
	CursorID        uint64
	Mode            uint8
	PublicInstaceID string
	IP              net.IP

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
	UTM             utm
	GeoResult       geoResult
	ScreenInfo      screenInfo
	BreadCrumb      breadCrumb
	Performance     performance

	PropsUserIDOrName string
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
	recordModeMap["pv_js"] = recordModePageViewJavaScript
	recordModeMap["pv_il"] = recordModePageViewImageLegacy
	recordModeMap["pv_ins"] = recordModePageViewImageNoScript
	recordModeMap["pv_amp"] = recordModePageViewAMP
	recordModeMap["pv_amp_i"] = recordModePageViewAMPImage

	recordModeMap["e_o"] = recordModeEventOther
	recordModeMap["e_js_pv"] = recordModeEventJSInPageView
	recordModeMap["e_js_c"] = recordModeEventJSCross
	recordModeMap["e_sw"] = recordModeEventServiceWorker
	recordModeMap["e_api"] = recordModeEventAPI

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

	r.Mode = mode
	r.PublicInstaceID = publicInstaceID

	return &r, nil
}

func (r *record) setAPI(
	postRequestAPI *postRequestAPI,
) *errorMessage {
	if postRequestAPI == nil {
		return &error_api_fields_missed
	}

	n := time.Now().Unix()
	if postRequestAPI.ClientTime <= n && postRequestAPI.ClientTime >= (n-28800) {
		r.Created = time.Unix(postRequestAPI.ClientTime, 0)
	}

	ip := net.ParseIP(postRequestAPI.ClientIP)
	if ip == nil {
		return &error_api_client_id_not_valid
	}

	r.IP = ip

	if postRequestAPI.ClientUserAgent == "" {
		return &error_api_client_user_agent_not_valid
	}

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

func (r *record) isClientError() bool {
	return r.Mode == recordModeClientError
}

func (r *record) verify(
	projectsManager *projects,
	privateKey string,
) *errorMessage {
	// initialize
	if r.Mode < 1 || r.IP == nil {
		return &error_record_not_proccessed
	}

	// in api mode private key must matched
	if r.Mode == recordModeEventAPI && !projectsManager.validateIDAndPrivate(r.PublicInstaceID, privateKey) {
		return &error_api_invalid_private_key
	}

	// in page js event must match with page url
	if r.Mode == recordModeEventJSInPageView && !projectsManager.validateIDAndURL(r.PublicInstaceID, r.PURL) {
		return &error_project_public_id_url_did_not_matched
	}

	// page view require matched project id and url of page view
	if r.isPageView() {
		if r.PURL == nil {
			return &error_url_required_and_must_valid
		}
		if !projectsManager.validateIDAndURL(r.PublicInstaceID, r.PURL) {
			return &error_project_public_id_url_did_not_matched
		}
		return nil
	}

	if r.isClientError() && r.PURL != nil && !projectsManager.validateIDAndURL(r.PublicInstaceID, r.PURL) {
		return &error_project_public_id_url_did_not_matched
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
	if postRequest.Page != nil {
		r.PURL = getURL(postRequest.Page.URL)
		r.PropsUserIDOrName = postRequest.Page.PropsUserIDOrName
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

		if postRequest.Page.RefererSessionURL != "" {
			r.SRefererURL = refererParser.parse(r.PURL, getURL(postRequest.Page.RefererSessionURL))
		}

		if postRequest.Page.PerformanceData != nil {
			r.Performance.PerfIsProcessed = true
			r.Performance.PerfPageLoadTime = postRequest.Page.PerformanceData.PerfPageLoadTime
			r.Performance.PerfDomainLookupTime = postRequest.Page.PerformanceData.PerfDomainLookupTime
			r.Performance.PerfTCPConnectTime = postRequest.Page.PerformanceData.PerfTCPConnectTime
			r.Performance.PerfServerResponseTime = postRequest.Page.PerformanceData.PerfServerResponseTime
			r.Performance.PerfPageDownloadTime = postRequest.Page.PerformanceData.PerfPageDownloadTime
			r.Performance.PerfRedirectTime = postRequest.Page.PerformanceData.PerfRedirectTime
			r.Performance.PerfDOMInteractiveTime = postRequest.Page.PerformanceData.PerfDOMInteractiveTime
			r.Performance.PerfContentLoadTime = postRequest.Page.PerformanceData.PerfContentLoadTime
			r.Performance.PerfResource = postRequest.Page.PerformanceData.PerfResource
		}

		if postRequest.Page.PageBreadcrumbObject != nil {
			r.BreadCrumb.BCIsProcessed = true
			r.BreadCrumb.BCN1 = postRequest.Page.PageBreadcrumbObject.N1
			r.BreadCrumb.BCN2 = postRequest.Page.PageBreadcrumbObject.N2
			r.BreadCrumb.BCN3 = postRequest.Page.PageBreadcrumbObject.N3
			r.BreadCrumb.BCN4 = postRequest.Page.PageBreadcrumbObject.N4
			r.BreadCrumb.BCN5 = postRequest.Page.PageBreadcrumbObject.N5
			r.BreadCrumb.BCP1 = postRequest.Page.PageBreadcrumbObject.P1
			r.BreadCrumb.BCP2 = postRequest.Page.PageBreadcrumbObject.P2
			r.BreadCrumb.BCP3 = postRequest.Page.PageBreadcrumbObject.P3
			r.BreadCrumb.BCP4 = postRequest.Page.PageBreadcrumbObject.P4
			r.BreadCrumb.BCP5 = postRequest.Page.PageBreadcrumbObject.P5
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
			re := recordEvent{
				ECategory: ev.Category,
				EAction:   ev.Action,
				ELabel:    ev.Label,
				EValue:    ev.Value,
			}
			events = append(events, re)
		}
		r.Events = events
		r.EventCount = len(events)
	}

	if r.Mode == recordModePageViewJavaScript && postRequest.CIDStd != "" {
		cid, cidErr := clientIDStandardParser(postRequest.CIDStd)
		if cidErr == nil {
			r.CID = cid
		}
	} else if r.Mode == recordModePageViewAMP && postRequest.CIDAmp != "" {
		r.CID = clientIDFromAMP(postRequest.CIDAmp)
	}

	if !r.CID.Valid && r.UserAgentResult.UaFull != "" && r.IP != nil {
		r.CID = clientIDFromOther([]string{r.IP.String(), r.UserAgentResult.UaFull})
	}
}

func (r *record) finalize() []byte {
	if r.PURL != nil {
		r.UTM = parseUTM(r.PURL)
	}

	if r.isPageView() {
		r.CursorID = getCursorID()
	}

	buf := &bytes.Buffer{}
	if err := gob.NewEncoder(buf).Encode(*r); err != nil {
		panic(err)
	}

	return buf.Bytes()
}
