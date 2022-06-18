package main

import (
	_ "embed"
	"strconv"

	ua "github.com/mileusna/useragent"
	"github.com/ua-parser/uap-go/uaparser"
)

// ineffassign: ignore
//go:embed embed/build/user_agents.yaml
var userAgents []byte

const (
	uaTypeUnknown = 0
	uaTypeMobile  = 1
	uaTypeDesktop = 2
	uaTypeTablet  = 3
	uaTypeBot     = 9
)

type userAgentParser struct {
	parser *uaparser.Parser
}

type userAgentResult struct {
	UaType                uint8
	UaFull                string
	UaChecksum            string
	UaBrowserName         string
	UaBrowserVersionMajor uint64
	UaBrowserVersion      string
	UaOSName              string
	UaOSVersionMajor      uint64
	UaOSVersion           string
	UaDeviceBrand         string
	UaDeviceFamily        string
	UaDeviceModel         string
}

func newUserAgentParser() *userAgentParser {
	uaParser := userAgentParser{}
	parser, err := uaparser.NewFromBytes(userAgents)
	if err != nil {
		panic(err)
	}
	uaParser.parser = parser
	return &uaParser
}

func (uaParser *userAgentParser) parse(uaString string) userAgentResult {
	client := uaParser.parser.Parse(uaString)

	result := userAgentResult{
		UaChecksum: checksum(uaString),
		UaFull:     uaString,
		UaType:     uaTypeUnknown,
	}

	// process type
	ua := ua.Parse(uaString)
	if ua.Mobile {
		result.UaType = uaTypeMobile
	} else if ua.Bot {
		result.UaType = uaTypeBot
	} else if ua.Desktop {
		result.UaType = uaTypeDesktop
	} else if ua.Tablet {
		result.UaType = uaTypeTablet
	}

	result.UaBrowserName = sanitizeTitle(client.UserAgent.Family)
	result.UaOSName = sanitizeTitle(client.Os.Family)
	result.UaDeviceBrand = sanitizeTitle(client.Device.Brand)
	result.UaDeviceFamily = sanitizeTitle(client.Device.Family)
	result.UaDeviceModel = sanitizeTitle(client.Device.Model)

	result.UaBrowserVersion = client.UserAgent.Major

	browserVersionMajor, browserVersionMajorErr := strconv.ParseUint(client.UserAgent.Major, 10, 64)
	if browserVersionMajorErr == nil && browserVersionMajor > 0 {
		result.UaBrowserVersionMajor = browserVersionMajor
	}

	if client.UserAgent.Minor != "" {
		result.UaBrowserVersion += "." + client.UserAgent.Minor
	}
	if client.UserAgent.Patch != "" {
		result.UaBrowserVersion += "." + client.UserAgent.Patch
	}

	result.UaOSVersion = client.Os.Major

	osVersionMajor, osVersionMajorErr := strconv.ParseUint(client.Os.Major, 10, 64)
	if osVersionMajorErr == nil && osVersionMajor > 0 {
		result.UaOSVersionMajor = osVersionMajor
	}

	if client.Os.Minor != "" {
		result.UaOSVersion += "." + client.Os.Minor
	}
	if client.Os.Patch != "" {
		result.UaOSVersion += "." + client.Os.Patch
	}

	return result
}
