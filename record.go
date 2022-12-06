package main

import (
	"encoding/json"
	"errors"
	"net"
	"net/url"
	"strings"
	"time"
)

type breadCrumb struct {
	BCIsProcessed bool   `json:"b"`
	BCLevel       uint8  `json:"l"`
	BCN1          string `json:"n1"`
	BCN2          string `json:"n2"`
	BCN3          string `json:"n3"`
	BCN4          string `json:"n4"`
	BCN5          string `json:"n5"`
	BCP1          string `json:"p1"`
	BCP2          string `json:"p2"`
	BCP3          string `json:"p3"`
	BCP4          string `json:"p4"`
	BCP5          string `json:"p5"`
}

type segment struct {
	S1N string `json:"n1"`
	S2N string `json:"n2"`
	S3N string `json:"n3"`
	S4N string `json:"n4"`
	S5N string `json:"n5"`

	S1V string `json:"v1"`
	S2V string `json:"v2"`
	S3V string `json:"v3"`
	S4V string `json:"v4"`
	S5V string `json:"v5"`
}

type performance struct {
	PerfIsProcessed        bool   `json:"p"`
	PerfPageLoadTime       uint16 `json:"pl"`
	PerfDomainLookupTime   uint16 `json:"dl"`
	PerfTCPConnectTime     uint16 `json:"t"`
	PerfServerResponseTime uint16 `json:"st"`
	PerfPageDownloadTime   uint16 `json:"pt"`
	PerfRedirectTime       uint16 `json:"rd"`
	PerfDOMInteractiveTime uint16 `json:"d"`
	PerfContentLoadTime    uint16 `json:"c"`
	PerfResource           uint16 `json:"r"`
}

type recordEvent struct {
	ECategory string `json:"c"`
	EAction   string `json:"a"`
	ELabel    string `json:"l"`
	EIdent    string `json:"i"`
	EValue    uint64 `json:"v"`
}

type record struct {
	ClientErrorMessage string    `json:"cem"`
	ClientErrorObject  string    `json:"ceo"`
	Created            time.Time `json:"c"`
	CursorID           uint64    `json:"cur"`
	Mode               uint8     `json:"m"`
	modeString         string
	PublicInstanceID   string `json:"p"`
	IP                 net.IP `json:"ip"`

	PEntityID         string `json:"ei"`
	PEntityModule     string `json:"em"`
	PEntityTaxonomyID uint16 `json:"et"`
	PURL              string `json:"u"`
	pURL              *url.URL
	PCanonicalURL     string `json:"cu"`
	pCanonicalURL     *url.URL
	PTitle            string      `json:"t"`
	PLang             string      `json:"l"`
	PIsIframe         bool        `json:"if"`
	PIsTouchSupport   bool        `json:"ts"`
	PKeywords         []string    `json:"k"`
	PRefererURL       refererData `json:"r"`
	SRefererURL       refererData `json:"sr"`

	CID clientID `json:"cid"`

	EventCount int           `json:"ec"`
	Events     []recordEvent `json:"e"`

	UserAgentResult userAgentResult `json:"ua"`
	Utm             utm             `json:"ut"`
	GeoResult       geoResult       `json:"g"`
	ScreenInfo      screenInfo      `json:"s"`
	BreadCrumb      breadCrumb      `json:"b"`
	Performance     performance     `json:"pr"`
	Segments        segment         `json:"sg"`
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
	recordModeClientErrorLegacy uint8 = 254
	recordModeClientError       uint8 = 255
)

var recordModeMap map[string]uint8

func init() {
	recordModeMap = make(map[string]uint8)

	// page view
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
	recordModeMap["err_l"] = recordModeClientErrorLegacy
}

func validateMode(m string) (uint8, error) {
	if value, ok := recordModeMap[m]; ok {
		return value, nil
	}
	return 0, errors.New("invalid mode of request")
}

func newRecord(modeQuery string, publicInstanceIDQuery string) (*record, error) {
	r := record{
		Created: time.Now(),
	}

	mode, modeErr := validateMode(modeQuery)

	if modeErr != nil {
		return nil, modeErr
	}

	publicInstanceID, publicInstanceIDErr := validatePublicInstanceID(publicInstanceIDQuery)

	if publicInstanceIDErr != nil {
		return nil, publicInstanceIDErr
	}

	r.modeString = modeQuery
	r.Mode = mode
	r.PublicInstanceID = publicInstanceID

	return &r, nil
}

func (r *record) setAPI(
	projectsManager *projects,
	userAgentParser *userAgentParser,
	geoParser *geoParser,
	postData *postRequest,

) *errorMessage {
	if postData.API == nil {
		return &errorAPIFieldsAreMissing
	}

	if r.isAPI() && !projectsManager.validateIDAndPrivate(r.PublicInstanceID, postData.API.PrivateInstanceKey) {
		return &errorAPIPrivateKeyFailed
	}

	n := time.Now().Unix()
	if postData.API.ClientTime <= n && postData.API.ClientTime >= (n-86400) {
		r.Created = time.Unix(postData.API.ClientTime, 0)
	}

	ip := net.ParseIP(postData.API.ClientIP).To4()
	if ip == nil {
		return &errorAPIClientIPNotValid
	}

	r.IP = ip

	if postData.API.ClientUserAgent == "" {
		return &errorAPIClientUserAgentNotValid
	}

	r.CID = clientIDNoneSTD([]string{r.IP.String(), postData.API.ClientUserAgent}, clientIDTypeOther)

	// apply updates
	r.UserAgentResult = userAgentParser.parse(postData.API.ClientUserAgent)
	r.GeoResult = geoParser.newResultFromIP(ip)

	return nil
}

func (r *record) isAPI() bool {
	return r.Mode == recordModeEventAPI
}

func (r *record) isPageView() bool {
	return r.Mode < 100
}

func (r *record) isEvent() bool {
	return r.Mode >= 100 && r.Mode < 200
}

func (r *record) isClientError() bool {
	return r.Mode >= 200
}

func (r *record) isImage() bool {
	if r.Mode == recordModePageViewAMPImage ||
		r.Mode == recordModePageViewImageLegacy ||
		r.Mode == recordModeClientErrorLegacy ||
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
		return &errorRecordNotProcessedYet
	}

	// page view require matched project id and url of page view
	if r.isPageView() {
		if r.pURL == nil {
			return &errorURLRequiredAndMustBeValid
		}
		if !projectsManager.validateIDAndURL(r.PublicInstanceID, r.pURL) {
			return &errorProjectPublicIDAndURLDidNotMatched
		}
	}

	if r.isAPI() && !projectsManager.validateIDAndPrivate(r.PublicInstanceID, privateKey) {
		return &errorAPIPrivateKeyFailed
	}

	if r.isClientError() && (r.PURL == "" || !projectsManager.validateIDAndURL(r.PublicInstanceID, r.pURL)) {
		return &errorProjectPublicIDAndURLDidNotMatched
	}

	if r.Mode == recordModeEventJSInPageView && !projectsManager.validateIDAndURL(r.PublicInstanceID, r.pURL) {
		return &errorProjectPublicIDAndURLDidNotMatched
	}

	if r.isEvent() && r.EventCount < 1 {
		return &errorEventsAreEmpty
	}

	return nil
}

func (r *record) setReferer(refererParser *refererParser, refererURL *url.URL) {
	r.PRefererURL = refererParser.parse(getURL(r.PURL), refererURL)
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
	r.PURL = sanitizeURL(qURL)
	r.PCanonicalURL = sanitizeURL(qCanonical)
	r.PTitle = sanitizeText(qTitle)
	r.PLang = sanitizeLanguage(qLang)
	r.PEntityID = sanitizeEntityID(qEntityID)
	r.PEntityModule = sanitizeName(qEntityModule)
	r.PEntityTaxonomyID = sanitizeEntityTaxonomyID(qEntityTaxonomyID)

	r.pURL = getURL(r.PURL)
	r.pCanonicalURL = getURL(r.PCanonicalURL)
}

func (r *record) setPostRequest(
	postRequest *postRequest,
	refererParser *refererParser,
	geoParser *geoParser,
) {
	if postRequest.ClientErrorMessage != "" {
		r.ClientErrorMessage = sanitizeText(postRequest.ClientErrorMessage)
		r.ClientErrorObject = sanitizeText(postRequest.ClientErrorObject)
	}

	if postRequest.Page != nil {
		r.PURL = sanitizeURL(postRequest.Page.URL)
		r.PCanonicalURL = sanitizeURL(postRequest.Page.CanonicalURL)
		r.PTitle = sanitizeText(postRequest.Page.Title)
		r.PLang = sanitizeLanguage(postRequest.Page.Lang)
		r.PEntityID = sanitizeEntityID(postRequest.Page.MainEntityID)
		r.PEntityModule = sanitizeName(postRequest.Page.MainEntityModule)
		r.PEntityTaxonomyID = sanitizeEntityTaxonomyID(postRequest.Page.MainEntityTaxonomyID)

		r.pURL = getURL(r.PURL)
		r.pCanonicalURL = getURL(r.PCanonicalURL)

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
			r.PRefererURL = refererParser.parse(r.pURL, getURL(postRequest.Page.RefererURL))
		}

		if postRequest.Page.RefererSessionURL != "" {
			r.SRefererURL = refererParser.parse(r.pURL, getURL(postRequest.Page.RefererSessionURL))
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

		if postRequest.Page.PageBreadcrumbObject != nil && r.pURL != nil {
			u1 := getURL(postRequest.Page.PageBreadcrumbObject.U1)
			pu := r.pURL
			if u1 != nil && u1.Host == pu.Host && postRequest.Page.PageBreadcrumbObject.N1 != "" {
				r.BreadCrumb.BCIsProcessed = true
				r.BreadCrumb.BCLevel = 1
				r.BreadCrumb.BCN1 = sanitizeText(postRequest.Page.PageBreadcrumbObject.N1)
				r.BreadCrumb.BCP1 = getURLPath(getURL(postRequest.Page.PageBreadcrumbObject.U1))

				u2 := getURL(postRequest.Page.PageBreadcrumbObject.U2)
				if u2 != nil && u2.Host == pu.Host && postRequest.Page.PageBreadcrumbObject.N2 != "" {
					r.BreadCrumb.BCLevel = 2
					r.BreadCrumb.BCN2 = sanitizeText(postRequest.Page.PageBreadcrumbObject.N2)
					r.BreadCrumb.BCP2 = getURLPath(getURL(postRequest.Page.PageBreadcrumbObject.U2))

					u3 := getURL(postRequest.Page.PageBreadcrumbObject.U3)
					if u3 != nil && u3.Host == pu.Host && postRequest.Page.PageBreadcrumbObject.N3 != "" {
						r.BreadCrumb.BCLevel = 3
						r.BreadCrumb.BCN3 = sanitizeText(postRequest.Page.PageBreadcrumbObject.N3)
						r.BreadCrumb.BCP3 = getURLPath(getURL(postRequest.Page.PageBreadcrumbObject.U3))

						u4 := getURL(postRequest.Page.PageBreadcrumbObject.U4)
						if u4 != nil && u4.Host == pu.Host && postRequest.Page.PageBreadcrumbObject.N4 != "" {
							r.BreadCrumb.BCLevel = 4
							r.BreadCrumb.BCN4 = sanitizeText(postRequest.Page.PageBreadcrumbObject.N4)
							r.BreadCrumb.BCP4 = getURLPath(getURL(postRequest.Page.PageBreadcrumbObject.U4))

							u5 := getURL(postRequest.Page.PageBreadcrumbObject.U5)
							if u5 != nil && u5.Host == pu.Host && postRequest.Page.PageBreadcrumbObject.N5 != "" {
								r.BreadCrumb.BCLevel = 5
								r.BreadCrumb.BCN5 = sanitizeText(postRequest.Page.PageBreadcrumbObject.N5)
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

		if postRequest.Page.Seg != nil {
			S1Name := sanitizeName(postRequest.Page.Seg.S1N)
			S1Value := strings.TrimSpace(postRequest.Page.Seg.S1V)
			if S1Name != "" && S1Value != "" {
				r.Segments.S1N = S1Name
				r.Segments.S1V = sanitizeText(S1Value)
			}

			S2Name := sanitizeName(postRequest.Page.Seg.S2N)
			S2Value := strings.TrimSpace(postRequest.Page.Seg.S2V)
			if S2Name != "" && S2Value != "" {
				r.Segments.S2N = S2Name
				r.Segments.S2V = sanitizeText(S2Value)
			}

			S3Name := sanitizeName(postRequest.Page.Seg.S3N)
			S3Value := strings.TrimSpace(postRequest.Page.Seg.S3V)
			if S3Name != "" && S3Value != "" {
				r.Segments.S3N = S3Name
				r.Segments.S3V = sanitizeText(S3Value)
			}

			S4Name := sanitizeName(postRequest.Page.Seg.S4N)
			S4Value := strings.TrimSpace(postRequest.Page.Seg.S4V)
			if S4Name != "" && S4Value != "" {
				r.Segments.S4N = S4Name
				r.Segments.S4V = sanitizeText(S4Value)
			}

			S5Name := sanitizeName(postRequest.Page.Seg.S5N)
			S5Value := strings.TrimSpace(postRequest.Page.Seg.S5V)
			if S5Name != "" && S5Value != "" {
				r.Segments.S5N = S5Name
				r.Segments.S5V = sanitizeText(S5Value)
			}
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
					EIdent:    sanitizeEntityID(ev.Ident),
					EValue:    uint64(ev.Value),
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
		r.CID = clientIDNoneSTD([]string{postRequest.CIDAmp}, clientIDTypeAmp)
	}
}

func (r *record) finalize() ([]byte, *errorMessage) {
	if r.Mode < 1 || !r.CID.Valid {
		return nil, &errorInvalidModeOrProjectPublicID
	}

	if r.pURL != nil {
		r.Utm = parseUTM(r.pURL)
	}

	if r.isPageView() {
		cursorID, cursorIDErr := getCursorID()
		if cursorIDErr != nil {
			e := errorInternalDependencyFailed
			e.debug = cursorIDErr.Error()
			return nil, &e
		}
		r.CursorID = cursorID
	}

	bytes, bytesErr := json.Marshal(r)

	if bytesErr != nil {
		e := errorInternalDependencyFailed
		e.debug = "marshal: " + bytesErr.Error()
		return nil, &e
	}

	return bytes, nil
}
