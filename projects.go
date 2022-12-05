package main

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"
)

type projects struct {
	sync.Mutex
	total             int
	publicIDs         map[string]bool
	privateKeys       map[string]string
	cached            map[string]string
	domainMap         map[string]string
	wildcardDomainMap map[string]string
}

type projectData struct {
	PrivateKey      string   `json:"p"`
	Domains         []string `json:"d,omitempty"`
	WildcardDomains []string `json:"w,omitempty"`
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
	total := 0

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
		total += 1
	}

	p.total = total
	p.cached = cached
	p.domainMap = domainMap
	p.wildcardDomainMap = wildcardDomainMap
	p.publicIDs = publicIDs
	p.privateKeys = privateKeys

	return nil
}

func (p *projects) validateIDAndPrivate(publicInstanceID string, privateKey string) bool {
	p.Lock()
	defer p.Unlock()

	if publicInstanceIDFromPrivate, ok := p.privateKeys[privateKey]; ok && privateKey != "" && publicInstanceIDFromPrivate == publicInstanceID {
		return true
	}

	return false
}

func (p *projects) validateID(publicInstanceID string) bool {
	p.Lock()
	defer p.Unlock()

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

func projectsLoadJSON(pathJSON string) (map[string]projectData, error) {
	b, err := os.ReadFile(pathJSON)
	if err != nil {
		return nil, err
	}
	var r map[string]projectData
	errJSON := json.Unmarshal(b, &r)
	if errJSON != nil {
		return nil, err
	}
	return r, nil
}

func projectsLoad(url string) (map[string]projectData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var r map[string]projectData
	errJSON := json.Unmarshal(body, &r)
	if errJSON != nil {
		return nil, err
	}
	return r, nil
}

func validatePublicInstanceIDRegex(publicInstanceID string) bool {
	return publicInstanceIDRegex.MatchString(publicInstanceID)
}
