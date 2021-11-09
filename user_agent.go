package main

import (
	_ "embed"
	"strconv"

	"github.com/ua-parser/uap-go/uaparser"
)

// ineffassign: ignore
//go:embed embed/build/user_agents.yaml
var userAgents []byte

// UserAgentParser is parser instance
type UserAgentParser struct {
	parser *uaparser.Parser
}

// UserAgentResult result of processed data
type UserAgentResult struct {
	BrowserName         string
	BrowserVersionMajor uint64
	BrowserVersion      string
	OSName              string
	OSVersionMajor      uint64
	OSVersion           string
	DeviceBrand         string
	DeviceFamily        string
	DeviceModel         string
}

// NewUserAgentParser new instance of UserAgentParser
func NewUserAgentParser() *UserAgentParser {
	uaParser := UserAgentParser{}
	parser, err := uaparser.NewFromBytes(userAgents)
	if err != nil {
		panic(err)
	}
	uaParser.parser = parser
	return &uaParser
}

// Parse return result for user agent
func (uaParser *UserAgentParser) Parse(uaString string) UserAgentResult {
	client := uaParser.parser.Parse(uaString)

	result := UserAgentResult{}

	result.BrowserName = client.UserAgent.Family
	result.OSName = client.Os.Family
	result.DeviceBrand = client.Device.Brand
	result.DeviceFamily = client.Device.Family
	result.DeviceModel = client.Device.Model

	result.BrowserVersion = client.UserAgent.Major

	browserVersionMajor, err1 := strconv.ParseUint(client.UserAgent.Major, 10, 64)
	if err1 == nil && browserVersionMajor > 0 {
		result.BrowserVersionMajor = browserVersionMajor
	}

	if client.UserAgent.Minor != "" {
		result.BrowserVersion += "." + client.UserAgent.Minor
	}
	if client.UserAgent.Patch != "" {
		result.BrowserVersion += "." + client.UserAgent.Patch
	}

	result.OSVersion = client.Os.Major

	osVersionMajor, err2 := strconv.ParseUint(client.Os.Major, 10, 64)
	if err2 == nil && osVersionMajor > 0 {
		result.OSVersionMajor = osVersionMajor
	}

	if client.Os.Minor != "" {
		result.OSVersion += "." + client.Os.Minor
	}
	if client.Os.Patch != "" {
		result.OSVersion += "." + client.Os.Patch
	}

	return result
}
