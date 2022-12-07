package main

import (
	_ "embed"
	"strconv"

	ua "github.com/mileusna/useragent"
	"github.com/ua-parser/uap-go/uaparser"
)

// ineffassign: ignore
//
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
	UaType                uint8  `json:"t"`
	UaFull                string `json:"f"`
	UaChecksum            string `json:"c"`
	UaBrowserName         string `json:"b"`
	UaBrowserVersionMajor uint64 `json:"bvm"`
	UaBrowserVersion      string `json:"bv"`
	UaOSName              string `json:"o"`
	UaOSVersionMajor      uint64 `json:"ovm"`
	UaOSVersion           string `json:"ov"`
	UaDeviceBrand         string `json:"db"`
	UaDeviceFamily        string `json:"df"`
	UaDeviceModel         string `json:"dm"`
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
		UaFull:     sanitizeText(uaString),
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
