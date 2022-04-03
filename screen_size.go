package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var sizeRegex = regexp.MustCompile(`^([0-9]{2,6})x([0-9]{2,6})$`)

type screenInfo struct {
	ScrIsProcessed                  bool
	ScrScreenOrientation            bool
	ScrScreenOrientationIsPortrait  bool
	ScrScreenOrientationIsSecondary bool
	ScrScreen                       string
	ScrScreenWidth                  uint16
	ScrScreenHeight                 uint16
	ScrViewport                     string
	ScrViewportWidth                uint16
	ScrViewportHeight               uint16
	ScrResoluton                    string
	ScrResolutonWidth               uint16
	ScrResolutonHeight              uint16
	ScrDevicePixelRatio             float64
	ScrColorDepth                   uint8
}

func parseScreenSize(
	screenOrientation string,
	screenSize string,
	viewportSize string,
	devicePixelRatio string,
	colorDepth string,
) screenInfo {
	result := screenInfo{
		ScrIsProcessed:       true,
		ScrScreenOrientation: false,
	}

	if devicePixelRatioFloat, err := strconv.ParseFloat(devicePixelRatio, 64); err == nil {
		result.ScrDevicePixelRatio = devicePixelRatioFloat
	}

	if colorDepthUInt8, err := strconv.ParseInt(colorDepth, 10, 16); err == nil {
		result.ScrColorDepth = uint8(colorDepthUInt8)
	}

	switch screenOrientation {
	case "p-p":
		result.ScrScreenOrientation = true
		result.ScrScreenOrientationIsPortrait = true
		result.ScrScreenOrientationIsSecondary = false
	case "p-s":
		result.ScrScreenOrientation = true
		result.ScrScreenOrientationIsPortrait = true
		result.ScrScreenOrientationIsSecondary = true
	case "l-p":
		result.ScrScreenOrientation = true
		result.ScrScreenOrientationIsPortrait = false
		result.ScrScreenOrientationIsSecondary = false
	case "l-s":
		result.ScrScreenOrientation = true
		result.ScrScreenOrientationIsPortrait = false
		result.ScrScreenOrientationIsSecondary = true
	}

	if sizeRegex.MatchString(viewportSize) {
		matched := sizeRegex.FindStringSubmatch(viewportSize)
		viewportWidth, viewportWidthErr := strconv.ParseUint(matched[1], 10, 16)
		viewportHeight, viewportHeightErr := strconv.ParseUint(matched[2], 10, 16)

		if viewportWidthErr == nil && viewportHeightErr == nil {
			result.ScrViewportWidth = uint16(viewportWidth)
			result.ScrViewportHeight = uint16(viewportHeight)
			result.ScrViewport = fmt.Sprintf("%dx%d", result.ScrViewportWidth, result.ScrViewportHeight)
		}
	}

	if sizeRegex.MatchString(screenSize) {
		matched := sizeRegex.FindStringSubmatch(screenSize)
		screenWidth, screenWidthErr := strconv.ParseUint(matched[1], 10, 16)
		screenHeight, screenHeightErr := strconv.ParseUint(matched[2], 10, 16)

		if screenWidthErr == nil && screenHeightErr == nil {
			result.ScrScreenWidth = uint16(screenWidth)
			result.ScrScreenHeight = uint16(screenHeight)
			result.ScrScreen = fmt.Sprintf("%dx%d", screenWidth, screenHeight)

			if result.ScrDevicePixelRatio != 0 {
				result.ScrResolutonWidth = uint16(float64(screenWidth) * result.ScrDevicePixelRatio)
				result.ScrResolutonHeight = uint16(float64(screenHeight) * result.ScrDevicePixelRatio)
				result.ScrResoluton = fmt.Sprintf("%dx%d", result.ScrResolutonWidth, result.ScrResolutonHeight)
			}
		}
	}

	return result
}
