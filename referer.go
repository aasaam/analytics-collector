package main

import (
	_ "embed"
	"net/url"

	"gopkg.in/yaml.v2"
)

//go:embed embed/referer.yml
var referer []byte

const (
	refererTypeNone         uint8 = 0
	refererTypeOther        uint8 = 1
	refererTypeSearchEngine uint8 = 2
	refererTypeSocial       uint8 = 3
	refererTypeMailProvider uint8 = 4
)

type refererParser struct {
	domainMap domainMap
}

type refererData struct {
	RefExist          bool   `json:"e"`
	RefURL            string `json:"u"`
	RefName           string `json:"n"`
	RefExternalHost   string `json:"h"`
	RefExternalDomain string `json:"d"`
	RefScheme         string `json:"s"`
	RefType           uint8  `json:"t"`
}

type domainMap map[string]*refererData

type refererYAML struct {
	SearchEngine map[string][]string `yaml:"search_engine"`
	SocialMedia  map[string][]string `yaml:"social_media"`
	MailProvider map[string][]string `yaml:"mail_provider"`
}

func getRefererData(name string, refererType uint8) *refererData {
	i := refererData{
		RefName: name,
		RefType: refererType,
	}

	return &i
}

func newRefererParser() *refererParser {
	domainMap := domainMap{}
	refererYAML := refererYAML{}
	err := yaml.Unmarshal(referer, &refererYAML)
	if err != nil {
		panic(err)
	}

	// SearchEngine
	for name, domains := range refererYAML.SearchEngine {
		for _, domain := range domains {
			domainMap[domain] = getRefererData(name, refererTypeSearchEngine)
		}
	}

	// SocialMedia
	for name, domains := range refererYAML.SocialMedia {
		for _, domain := range domains {
			domainMap[domain] = getRefererData(name, refererTypeSocial)
		}
	}

	// MailProvider
	for name, domains := range refererYAML.MailProvider {
		for _, domain := range domains {
			domainMap[domain] = getRefererData(name, refererTypeMailProvider)
		}
	}

	result := refererParser{}
	result.domainMap = domainMap

	return &result
}

func (rp *refererParser) parse(currentURL *url.URL, refererURL *url.URL) refererData {
	result := refererData{
		RefType: refererTypeNone,
	}

	if refererURL == nil || currentURL == nil {
		return result
	}

	result.RefURL = refererURL.String()
	result.RefScheme = refererURL.Scheme

	rHost := getDomain(refererURL)
	cHost := getDomain(currentURL)

	if cHost != rHost {
		result.RefExist = true

		rDomain := getSecondDomainLevel(refererURL)
		cDomain := getSecondDomainLevel(currentURL)

		if rDomain != cDomain {
			result.RefExternalDomain = rDomain
		}

		result.RefExternalHost = rHost
		result.RefType = refererTypeOther
		result.RefName = rDomain
		if val, ok := rp.domainMap[rHost]; ok {
			result.RefName = val.RefName
			result.RefType = val.RefType
			return result
		}
	}

	return result
}
