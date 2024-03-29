package main

import (
	"net/url"
	"strings"
)

type utm struct {
	UtmValid    bool   `json:"v"`
	UtmExist    bool   `json:"e"`
	UtmSource   string `json:"s"`
	UtmMedium   string `json:"m"`
	UtmCampaign string `json:"c"`
	UtmID       string `json:"i"`
	UtmTerm     string `json:"t"`
	UtmContent  string `json:"cn"`
}

func parseUTM(u *url.URL) utm {
	result := utm{
		UtmExist: false,
		UtmValid: false,
	}

	if u == nil {
		return result
	}

	queries := u.Query()
	for k, v := range queries {
		kci := strings.ToLower(k)

		switch kci {
		case "utm_source":
			result.UtmExist = true
			result.UtmSource = v[0]
		case "utm_medium":
			result.UtmExist = true
			result.UtmMedium = v[0]
		case "utm_campaign":
			result.UtmExist = true
			result.UtmCampaign = v[0]
		case "utm_id":
			result.UtmExist = true
			result.UtmID = v[0]
		case "utm_term":
			result.UtmExist = true
			result.UtmTerm = v[0]
		case "utm_content":
			result.UtmExist = true
			result.UtmContent = v[0]
		}
	}

	result.UtmValid = false
	if result.UtmSource != "" && result.UtmMedium != "" && result.UtmCampaign != "" {
		result.UtmValid = true
	}

	return result
}
