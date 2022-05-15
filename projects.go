package main

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"sync"
)

type projects struct {
	sync.Mutex
	publicIDs         map[string]bool
	privateKeys       map[string]string
	cached            map[string]string
	domainMap         map[string]string
	wildcardDomainMap map[string]string
}

type projectData struct {
	PrivateKey      string   `json:"p"`
	Domains         []string `json:"d"`
	WildcardDomains []string `json:"w"`
}

var publicInstanceIDRegex = regexp.MustCompile(`^[a-zA-Z0-9]{12}$`)

func validatePublicInstanceID(pid string) (string, error) {
	if ok := publicInstanceIDRegex.MatchString(pid); ok {
		return pid, nil
	}
	return "", errors.New("invalid public instance id")
}

func newProjectsManager() *projects {
	cached := make(map[string]string)
	domainMap := make(map[string]string)
	wildcardDomainMap := make(map[string]string)
	privateKeys := make(map[string]string)
	publicIDs := make(map[string]bool)

	result := projects{
		cached:            cached,
		domainMap:         domainMap,
		wildcardDomainMap: wildcardDomainMap,
		privateKeys:       privateKeys,
		publicIDs:         publicIDs,
	}

	return &result
}

func (p *projects) load(data map[string]projectData) error {
	p.Lock()
	defer p.Unlock()

	cached := make(map[string]string)
	domainMap := make(map[string]string)
	wildcardDomainMap := make(map[string]string)
	privateKeys := make(map[string]string)
	publicIDs := make(map[string]bool)

	for publicID, projectData := range data {
		publicIDs[publicID] = true
		privateKeys[projectData.PrivateKey] = publicID
		for _, domain := range projectData.Domains {
			if isValidURL("http://" + domain) {
				domainMap[domain] = publicID
			}
		}
		for _, wildcardDomain := range projectData.WildcardDomains {
			if isValidURL("http://" + wildcardDomain) {
				wildcardDomainMap[wildcardDomain] = publicID
			}
		}
	}

	p.cached = cached
	p.domainMap = domainMap
	p.wildcardDomainMap = wildcardDomainMap
	p.publicIDs = publicIDs
	p.privateKeys = privateKeys

	return nil
}

func (p *projects) validateIDAndPrivate(publicInstanceID string, privateKey string) bool {
	if publicInstanceIDFromPrivate, ok := p.privateKeys[privateKey]; ok && privateKey != "" && publicInstanceIDFromPrivate == publicInstanceID {
		return true
	}

	return false
}

func (p *projects) validateID(publicInstanceID string) bool {
	return p.publicIDs[publicInstanceID]
}

func (p *projects) validateIDAndURL(publicInstanceID string, requestURL *url.URL) bool {
	if !publicInstanceIDRegex.MatchString(publicInstanceID) {
		return false
	}

	if requestURL == nil {
		return false
	}

	hostname := requestURL.Hostname()

	p.Lock()
	defer p.Unlock()

	if cached, ok := p.cached[publicInstanceID]; ok && cached == hostname {
		return true
	}

	if found, ok := p.domainMap[hostname]; ok {
		if found == publicInstanceID {
			p.cached[publicInstanceID] = hostname
			return true
		}
	}

	if net.ParseIP(hostname) != nil {
		return false
	}

	for domain, domainPublicInstanceID := range p.wildcardDomainMap {
		if domain == hostname {
			p.cached[publicInstanceID] = hostname
			return true
		}
		domainRegexString := `.*\.` + regexp.QuoteMeta(domain) + `$`
		domainRegex, domainRegexErr := regexp.Compile(domainRegexString)

		if domainRegexErr == nil && domainRegex.MatchString(hostname) && publicInstanceID == domainPublicInstanceID {
			p.cached[publicInstanceID] = hostname
			return true
		}
	}

	return false
}
