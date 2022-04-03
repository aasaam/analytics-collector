package main

import (
	"net/url"
	"strings"
)

type utm struct {
	UTMValid    bool
	UTMExist    bool
	UTMSource   string
	UTMMedium   string
	UTMCampaign string
	UTMID       string
	UTMTerm     string
	UTMContent  string
}

func parseUTM(u *url.URL) utm {
	result := utm{
		UTMExist: false,
		UTMValid: false,
	}

	if u == nil {
		return result
	}

	queries := u.Query()
	for k, v := range queries {
		kci := strings.ToLower(k)

		switch kci {
		case "utm_source":
			result.UTMExist = true
			result.UTMSource = v[0]
		case "utm_medium":
			result.UTMExist = true
			result.UTMMedium = v[0]
		case "utm_campaign":
			result.UTMExist = true
			result.UTMCampaign = v[0]
		case "utm_id":
			result.UTMExist = true
			result.UTMID = v[0]
		case "utm_term":
			result.UTMExist = true
			result.UTMTerm = v[0]
		case "utm_content":
			result.UTMExist = true
			result.UTMContent = v[0]
		}
	}

	result.UTMValid = false
	if result.UTMSource != "" && result.UTMMedium != "" && result.UTMCampaign != "" {
		result.UTMValid = true
	}

	return result
}
