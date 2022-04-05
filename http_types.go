package main

import (
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
)

const (
	nginxXAccelExpires      = "X-Accel-Expires"
	collectorURLReplacement = "__COLLECTOR_URL__"
)

const (
	errorTypeApp    = "app"
	errorTypeClient = "client"
)

const (
	mimetypeJS   = "text/javascript"
	mimetypeText = "text/plain"
	mimetypeGIF  = "image/gif"
)

const (
	metricsPath = "/metrics"
)

const (
	recordQueryMode            = "m"
	recordQueryPublicInstaceID = "i"

	recordQueryURL              = "u"
	recordQueryCanonicalURL     = "cu"
	recordQueryRefererURL       = "r"
	recordQueryTitle            = "t"
	recordQueryLang             = "l"
	recordQueryEntityID         = "ei"
	recordQueryEntityModule     = "em"
	recordQueryEntityTaxonomyID = "et"
)

type errorMessage struct {
	code int
	msg  string
}

type postRequestEvent struct {
	Category string `json:"c,omitempty"`
	Action   string `json:"a,omitempty"`
	Label    string `json:"l,omitempty"`
	Value    uint64 `json:"v,omitempty"`
}

type postRequestAPI struct {
	PrivateInstanceKey string `json:"i_p,omitempty"`
	ClientIP           string `json:"c_ip,omitempty"`
	ClientUserAgent    string `json:"c_ua,omitempty"`
	ClientTime         int64  `json:"c_t,omitempty"`
}

type postRequestGeographyData struct {
	Lat float64 `json:"lat,omitempty"`
	Lon float64 `json:"lon,omitempty"`
}

type postRequestBreadcrumb struct {
	N1 string `json:"n1,omitempty"`
	N2 string `json:"n2,omitempty"`
	N3 string `json:"n3,omitempty"`
	N4 string `json:"n4,omitempty"`
	N5 string `json:"n5,omitempty"`
	U1 string `json:"u1,omitempty"`
	U2 string `json:"u2,omitempty"`
	U3 string `json:"u3,omitempty"`
	U4 string `json:"u4,omitempty"`
	U5 string `json:"u5,omitempty"`
}

type postRequestPerformanceData struct {
	PerfPageLoadTime       string `json:"plt,omitempty"`
	PerfDomainLookupTime   string `json:"dlt,omitempty"`
	PerfTCPConnectTime     string `json:"tct,omitempty"`
	PerfServerResponseTime string `json:"srt,omitempty"`
	PerfPageDownloadTime   string `json:"pdt,omitempty"`
	PerfRedirectTime       string `json:"rt,omitempty"`
	PerfDOMInteractiveTime string `json:"dit,omitempty"`
	PerfContentLoadTime    string `json:"clt,omitempty"`
	PerfResource           uint16 `json:"r,omitempty"`
}

type postRequestPage struct {
	URL          string `json:"u,omitempty"`
	CanonicalURL string `json:"cu,omitempty"`
	RefererURL   string `json:"r,omitempty"`

	Title                string `json:"t,omitempty"`
	Lang                 string `json:"l,omitempty"`
	MainEntityID         string `json:"ei,omitempty"`
	MainEntityModule     string `json:"em,omitempty"`
	MainEntityTaxonomyID string `json:"et,omitempty"`

	ScreenSize            string `json:"scr,omitempty"`
	ViewportSize          string `json:"vps,omitempty"`
	ColorDepth            string `json:"cd,omitempty"`
	DevicePixelRatio      string `json:"dpr,omitempty"`
	ScreenOrientationType string `json:"sot,omitempty"`
	PageKeywords          string `json:"k,omitempty"`
	IsIframe              bool   `json:"if,omitempty"`
	IsTouchSupport        bool   `json:"ts,omitempty"`
	RefererSessionURL     string `json:"rs,omitempty"`

	PageBreadcrumbObject *postRequestBreadcrumb      `json:"bc,omitempty"`
	PerformanceData      *postRequestPerformanceData `json:"prf,omitempty"`
	GeographyData        *postRequestGeographyData   `json:"geo,omitempty"`

	UserIDOrName string `json:"usr,omitempty"`
}

type postRequest struct {
	ClientErrorMessage string              `json:"msg,omitempty"`
	ClientErrorObject  string              `json:"err,omitempty"`
	CIDStd             string              `json:"cid_std,omitempty"`
	CIDAmp             string              `json:"cid_amp,omitempty"`
	Page               *postRequestPage    `json:"p,omitempty"`
	Events             *[]postRequestEvent `json:"ev,omitempty"`
	API                *postRequestAPI     `json:"ar,omitempty"`
}

var singleGifImage, _ = base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")

var (
	errorInternalServerError                errorMessage = errorMessage{code: fiber.StatusInternalServerError, msg: "errorInternalServerError"}
	errorMetricsForbidden                   errorMessage = errorMessage{code: fiber.StatusForbidden, msg: "errorMetricsForbidden"}
	errorRecordNotProccessedYet             errorMessage = errorMessage{code: fiber.StatusConflict, msg: "errorRecordNotProccessedYet"}
	errorRecordNotValid                     errorMessage = errorMessage{code: fiber.StatusBadRequest, msg: "errorRecordNotValid"}
	errorRecordCIDNotProccessed             errorMessage = errorMessage{code: fiber.StatusBadRequest, msg: "errorRecordCIDNotProccessed"}
	errorProjectPublicIDAndURLDidNotMatched errorMessage = errorMessage{code: fiber.StatusForbidden, msg: "errorProjectPublicIDAndURLDidNotMatched"}
	errorInvalidModeOrProjectPublicID       errorMessage = errorMessage{code: fiber.StatusUnprocessableEntity, msg: "errorInvalidModeOrProjectPublicID"}
	errorURLRequiredAndMustBeValid          errorMessage = errorMessage{code: fiber.StatusFailedDependency, msg: "errorURLRequiredAndMustBeValid"}
	errorAPIFieldsAreMissing                errorMessage = errorMessage{code: fiber.StatusFailedDependency, msg: "errorAPIFieldsAreMissing"}
	errorAPIPrivateKeyFailed                errorMessage = errorMessage{code: fiber.StatusForbidden, msg: "errorAPIPrivateKeyFailed"}
	errorAPIClientIPNotValid                errorMessage = errorMessage{code: fiber.StatusForbidden, msg: "errorAPIClientIPNotValid"}
	errorAPIClientUserAgentNotValid         errorMessage = errorMessage{code: fiber.StatusForbidden, msg: "errorAPIClientUserAgentNotValid"}
	errorEventsAreEmpty                     errorMessage = errorMessage{code: fiber.StatusFailedDependency, msg: "errorEventsAreEmpty"}
	errorBadPOSTBody                        errorMessage = errorMessage{code: fiber.StatusUnprocessableEntity, msg: "errorBadPOSTBody"}
)