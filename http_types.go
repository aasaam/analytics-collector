package main

import (
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
)

const (
	nginx_x_accel_expires     = "X-Accel-Expires"
	collector_url_replacement = "__COLLECTOR_URL__"
)

const (
	error_type_app    = "app"
	error_type_client = "client"
)

const (
	mimetype_js  = "text/javascript"
	mimetype_map = "application/json"
	mimetype_gif = "image/gif"
)

const (
	metrics_path = "/metrics"
)

const (
	record_query_mode              = "m"
	record_query_public_instace_id = "i"

	record_query_url                = "u"
	record_query_canonical          = "cu"
	record_query_referer            = "r"
	record_query_title              = "t"
	record_query_lang               = "l"
	record_query_entity_id          = "ei"
	record_query_entity_module      = "em"
	record_query_entity_taxonomy_id = "et"
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

var single_gif_image, _ = base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")

var (
	error_record_not_proccessed                 errorMessage = errorMessage{code: fiber.StatusConflict, msg: "error_record_not_proccessed"}
	error_record_not_valid                      errorMessage = errorMessage{code: fiber.StatusBadRequest, msg: "error_record_not_valid"}
	error_record_cid_not_proccessed             errorMessage = errorMessage{code: fiber.StatusBadRequest, msg: "error_record_cid_not_proccessed"}
	error_project_public_id_url_did_not_matched errorMessage = errorMessage{code: fiber.StatusForbidden, msg: "error_project_public_id_url_did_not_matched"}
	error_invalid_mode_or_project_public_id     errorMessage = errorMessage{code: fiber.StatusUnprocessableEntity, msg: "error_invalid_mode_or_project_public_id"}
	error_url_required_and_must_valid           errorMessage = errorMessage{code: fiber.StatusFailedDependency, msg: "error_url_required_and_must_valid"}
	error_api_fields_missed                     errorMessage = errorMessage{code: fiber.StatusFailedDependency, msg: "error_api_fields_missed"}
	error_api_invalid_private_key               errorMessage = errorMessage{code: fiber.StatusForbidden, msg: "error_api_invalid_private_key"}
	error_api_client_id_not_valid               errorMessage = errorMessage{code: fiber.StatusForbidden, msg: "error_api_client_id_not_valid"}
	error_api_client_user_agent_not_valid       errorMessage = errorMessage{code: fiber.StatusForbidden, msg: "error_api_client_user_agent_not_valid"}
	error_events_are_empty                      errorMessage = errorMessage{code: fiber.StatusFailedDependency, msg: "error_events_are_empty"}
)
