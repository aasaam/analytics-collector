package main

import (
	"encoding/json"
	"net/url"
	"regexp"
	"sync"
)

type Projects struct {
	sync.RWMutex
	publicHash        map[string]bool
	privateHash       map[string]string
	cached            map[string]string
	domainMap         map[string]string
	wildcardDomainMap map[string]string
}

var projectIDRegex, _ = regexp.Compile(`^[a-zA-Z0-9]{12}$`)

var wildcardDomainRegex, _ = regexp.Compile(`^\*\.(.*)$`)

// NewProjectManager return project manager
func NewProjectManager() *Projects {
	cached := make(map[string]string)
	domainMap := make(map[string]string)
	wildcardDomainMap := make(map[string]string)
	privateHash := make(map[string]string)
	publicHash := make(map[string]bool)

	result := Projects{
		cached:            cached,
		domainMap:         domainMap,
		wildcardDomainMap: wildcardDomainMap,
		privateHash:       privateHash,
		publicHash:        publicHash,
	}

	return &result
}

// LoadJSON load json will store updated data for new requests
func (p *Projects) LoadJSON(jsonData []byte) error {
	type projectData struct {
		PrivateHash string   `json:"ph"`
		Domains     []string `json:"d"`
	}
	var data map[string]projectData

	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}

	p.Lock()
	defer p.Unlock()

	cached := make(map[string]string)
	privateHash := make(map[string]string)
	publicHash := make(map[string]bool)
	domainMap := make(map[string]string)
	wildcardDomainMap := make(map[string]string)

	for projectPublicHash, projectData := range data {
		publicHash[projectPublicHash] = true
		privateHash[projectData.PrivateHash] = projectPublicHash
		for _, domain := range projectData.Domains {
			// is wildcard?
			if wildcardDomainRegex.MatchString(domain) {
				matched := wildcardDomainRegex.FindStringSubmatch(domain)
				urlParsed, err := url.ParseRequestURI("http://" + matched[1])
				if err == nil {
					wildcardDomainMap[urlParsed.Hostname()] = projectPublicHash
				}
			} else { // normal domain
				urlParsed, err := url.ParseRequestURI("http://" + domain)
				if err == nil {
					domainMap[urlParsed.Hostname()] = projectPublicHash
				}
			}
		}
	}

	p.cached = cached
	p.domainMap = domainMap
	p.wildcardDomainMap = wildcardDomainMap
	p.privateHash = privateHash
	p.publicHash = publicHash

	return nil
}

// ValidateAPI is domain public hash validator
func (p *Projects) ValidateAPI(projectPublicHash string, projectPrivateHash string) bool {
	if publicHash, ok := p.privateHash[projectPrivateHash]; ok && publicHash == projectPublicHash {
		return true
	}

	return false
}

// ValidateEvent is check project public hash exist on request
func (p *Projects) ValidateEvent(projectPublicHash string) bool {
	return p.publicHash[projectPublicHash]
}

// Validate is domain public hash validator
func (p *Projects) ValidatePageView(projectPublicHash string, domain string) bool {
	if !projectIDRegex.MatchString(projectPublicHash) {
		return false
	}

	if cached, ok := p.cached[projectPublicHash]; ok && cached == domain {
		return true
	}

	p.RLock()
	defer p.RUnlock()

	if found, ok := p.domainMap[domain]; ok {
		if found == projectPublicHash {
			p.cached[projectPublicHash] = domain
			return true
		}
	}

	for mainHost, projectPublicHashIteration := range p.wildcardDomainMap {
		wildRegex1 := `^.*\.` + regexp.QuoteMeta(mainHost) + `$`

		regex1, e := regexp.Compile(wildRegex1)
		if e == nil && regex1.MatchString(domain) && projectPublicHash == projectPublicHashIteration {
			p.cached[projectPublicHash] = domain
			return true
		}
		wildRegex2 := `^` + regexp.QuoteMeta(mainHost) + `$`

		regex2, e := regexp.Compile(wildRegex2)
		if e == nil && regex2.MatchString(domain) && projectPublicHash == projectPublicHashIteration {
			p.cached[projectPublicHash] = domain
			return true
		}
	}

	return false
}
