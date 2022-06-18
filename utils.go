package main

import (
	"crypto/sha1"
	"encoding/base64"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/text/language"
)

const (
	checksumEmpty = "000000000000000000000000"
)

var sanitizeTitleRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)
var sanitizeTitleMoreSpaceRegex = regexp.MustCompile(`[\s]+`)
var sanitizeNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{1,31}$`)
var entityIDRegex = regexp.MustCompile(`^[a-zA-Z0-9-_\/]{1,63}$`)
var entityTaxonomyIDRegex = regexp.MustCompile(`^[A-Z]{1}[0-9a-z]{4}$`)
var checksumReplaceRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)
var cursorTimeLayout = "20060102030405.000"

func uint16FromString(s string) uint16 {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return uint16(i)
}

func boolUint8(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}

func getURLString(u *url.URL) string {
	if u == nil {
		return ""
	}
	return u.String()

}

func isValidURL(urlString string) bool {
	if urlString == "" {
		return false
	}

	_, err := url.Parse(urlString)

	return err == nil
}

func getURLPath(u *url.URL) string {
	r, e := regexp.Compile(regexp.QuoteMeta(u.Host) + `(.*)$`)
	if e != nil {
		return ""
	}
	if r.MatchString(u.String()) {
		matched := r.FindStringSubmatch(u.String())
		return matched[1]
	}
	return ""
}

func getURL(urlString string) *url.URL {
	if urlString == "" {
		return nil
	}

	u, err := url.Parse(urlString)
	if err == nil {
		return u
	}

	return nil
}

func sanitizeURL(urlString string) string {
	if urlString == "" {
		return ""
	}

	u, err := url.Parse(urlString)
	if err == nil {
		return u.String()
	}

	return ""
}

func checksum(str string) string {
	if strings.TrimSpace(str) == "" {
		return checksumEmpty
	}
	h := sha1.New()
	h.Write([]byte(str))
	return checksumReplaceRegex.ReplaceAllString(base64.StdEncoding.EncodeToString(h.Sum(nil)), "0")[0:24]
}

func sanitizeTitle(str string) string {
	s := sanitizeTitleRegex.ReplaceAllString(str, " ")
	s = sanitizeTitleMoreSpaceRegex.ReplaceAllString(s, " ")
	return strings.ToLower(strings.TrimSpace(s))
}

func normalizeHostname(hostname string) string {
	_, isICN := publicsuffix.PublicSuffix(hostname)
	if isICN {
		p := idna.New(idna.ValidateForRegistration())
		encodedHostname, err := p.ToUnicode(hostname)
		if err == nil {
			return encodedHostname
		}
	}
	return hostname
}

func getSecondDomainLevel(u *url.URL) string {
	hostname := u.Hostname()

	ip := net.ParseIP(hostname)
	if ip != nil {
		return ip.String()
	}

	hostname = normalizeHostname(hostname)

	parts := strings.Split(hostname, ".")

	if len(parts) >= 2 {
		tld := parts[len(parts)-1]
		domain := parts[len(parts)-2]
		return domain + "." + tld
	}

	return hostname
}

func getDomain(u *url.URL) string {
	hostname := u.Hostname()

	ip := net.ParseIP(hostname)
	if ip != nil {
		return ip.String()
	}

	return normalizeHostname(hostname)
}

func sanitizeName(name string) string {
	if ok := sanitizeNameRegex.MatchString(name); ok {
		return strings.ToLower(name)
	}
	return ""
}

func sanitizeEntityID(id string) string {
	if ok := entityIDRegex.MatchString(id); ok {
		return id
	}
	return ""
}

func sanitizeLanguage(locale string) string {
	tag, err := language.Parse(locale)
	if err != nil {
		return ""
	}
	base, _ := tag.Base()
	return base.String()
}

func sanitizeEntityTaxonomyID(id string) string {
	if ok := entityTaxonomyIDRegex.MatchString(id); ok {
		return id
	}
	return ""
}

func parseKeywords(inpKeywords string) []string {
	ks := strings.Split(inpKeywords, ",")
	r := []string{}
	for _, k := range ks {
		v := strings.TrimSpace(k)
		if len(v) >= 1 {
			r = append(r, v)
		}
	}

	if len(r) > 10 {
		return r[0:10]
	}

	return r
}

func getCursorID() uint64 {
	n := time.Now().Format(cursorTimeLayout)
	n = strings.ReplaceAll(n, ".", "")
	ui, err := strconv.ParseUint(n, 10, 64)
	if err != nil {
		panic(err)
	}
	return ui
}
