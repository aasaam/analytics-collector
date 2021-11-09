package main

import (
	_ "embed"
	"net/url"

	"gopkg.in/yaml.v2"
)

//go:embed embed/referrer.yml
var referrer []byte

const (
	// ReferrerTypeSearchEngine is referrer for search engines
	ReferrerTypeSearchEngine uint8 = 1
	// ReferrerTypeSocial is referrer for social media
	ReferrerTypeSocial uint8 = 2
	// ReferrerTypeMailProvider is referrer for mail providers
	ReferrerTypeMailProvider uint8 = 3
)

// ReferrerParser instace for domain map
type ReferrerParser struct {
	domainMap DomainMap
}

// RefererData is name and type of referrer
type RefererData struct {
	Name         string
	ExternalHost string
	ReferrerType uint8
}

// DomainMap map from string to DomainData
type DomainMap map[string]RefererData

// ReferrerYAML is yaml file structure for map data
type ReferrerYAML struct {
	SearchEngine map[string][]string `yaml:"search_engine"`
	SocialMedia  map[string][]string `yaml:"social_media"`
	MailProvider map[string][]string `yaml:"mail_provider"`
}

// Parse and return domain data if is special
func (rp *ReferrerParser) Parse(pageURL string, referrerURI string) RefererData {
	result := RefererData{}
	refU, err := url.Parse(referrerURI)
	if err != nil {
		return result
	}

	pageU, err := url.Parse(pageURL)
	if err != nil {
		return result
	}

	if pageU.Hostname() != refU.Hostname() {
		result.ExternalHost = refU.Hostname()
	}
	if val, ok := rp.domainMap[refU.Hostname()]; ok {
		result.Name = val.Name
		result.ReferrerType = val.ReferrerType
		return result
	}

	return result
}

// GetURLHost return hostname of url
func GetURLHost(urlString string) string {
	u, err := url.Parse(urlString)
	if err != nil {
		return ""
	}

	return u.Hostname()
}

// GetRefererData domain name to type of referrer
func GetRefererData(name string, referrerType uint8) RefererData {
	refererData := RefererData{
		Name:         name,
		ReferrerType: referrerType,
	}

	return refererData
}

// NewReferrerParser return instance of referrer parser
func NewReferrerParser() *ReferrerParser {
	domainMap := DomainMap{}
	referrerYAML := ReferrerYAML{}
	err := yaml.Unmarshal(referrer, &referrerYAML)
	if err != nil {
		panic(err)
	}

	// SearchEngine
	for name, domains := range referrerYAML.SearchEngine {
		for _, domain := range domains {
			domainMap[domain] = GetRefererData(name, ReferrerTypeSearchEngine)
		}
	}

	// SocialMedia
	for name, domains := range referrerYAML.SocialMedia {
		for _, domain := range domains {
			domainMap[domain] = GetRefererData(name, ReferrerTypeSocial)
		}
	}

	// MailProvider
	for name, domains := range referrerYAML.MailProvider {
		for _, domain := range domains {
			domainMap[domain] = GetRefererData(name, ReferrerTypeMailProvider)
		}
	}

	result := ReferrerParser{}
	result.domainMap = domainMap

	return &result
}
