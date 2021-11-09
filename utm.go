package main

import (
	"net/url"
	"strings"
)

// UTM is Urchin Tracking Module for online marketing campaigns
type UTM struct {
	Valid        bool
	Source       string
	Medium       string
	CampaignName string
	ID           string
	Term         string
	Content      string
}

// ParseUTM return parsed utm parameters from URL
func ParseUTM(urlString string) UTM {
	result := UTM{}

	u, err := url.Parse(urlString)
	if err != nil {
		return result
	}

	queries := u.Query()
	for k, v := range queries {
		kci := strings.ToLower(k)

		switch kci {
		case "utm_source":
			result.Source = strings.TrimSpace(v[0])
		case "utm_medium":
			result.Medium = strings.TrimSpace(v[0])
		case "utm_campaign":
			result.CampaignName = strings.TrimSpace(v[0])
		case "utm_id":
			result.ID = strings.TrimSpace(v[0])
		case "utm_term":
			result.Term = strings.TrimSpace(v[0])
		case "utm_content":
			result.Content = strings.TrimSpace(v[0])
		}
	}

	result.Valid = false
	if result.Source != "" && result.Medium != "" && result.CampaignName != "" {
		result.Valid = true
	}

	return result
}
