package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var sizeRegex, _ = regexp.Compile("^([0-9]{2,5})x([0-9]{2,5})$")

const (
	ScreenOrientationUnknown            uint8 = 0
	ScreenOrientationPortraitPrimary    uint8 = 1
	ScreenOrientationPortraitSecondary  uint8 = 2
	ScreenOrientationLandscapePrimary   uint8 = 3
	ScreenOrientationLandscapeSecondary uint8 = 4
)

// ScreenInfo is processed screen,viewport and resolution sizes
type ScreenInfo struct {
	ScreenOrientation uint8
	Screen            string
	ScreenWidth       uint16
	ScreenHeight      uint16
	Viewport          string
	ViewportWidth     uint16
	ViewportHeight    uint16
	Resoluton         string
	ResolutonWidth    uint16
	ResolutonHeight   uint16
	PixelRatio        float64
	ColorDepth        uint8
}

// ParseScreenSize return *ScreenInfo
func ParseScreenSize(screenOrientation string, screenSize string, viewportSize string, pixelRatio float64, colorDepth uint8) *ScreenInfo {
	result := ScreenInfo{
		ColorDepth:        colorDepth,
		PixelRatio:        pixelRatio,
		ScreenOrientation: ScreenOrientationUnknown,
	}

	switch screenOrientation {
	case "p-p":
		result.ScreenOrientation = ScreenOrientationPortraitPrimary
	case "p-s":
		result.ScreenOrientation = ScreenOrientationPortraitSecondary
	case "l-p":
		result.ScreenOrientation = ScreenOrientationLandscapePrimary
	case "l-s":
		result.ScreenOrientation = ScreenOrientationLandscapeSecondary
	}

	if sizeRegex.MatchString(viewportSize) {
		matched := sizeRegex.FindStringSubmatch(viewportSize)
		viewportWidth, _ := strconv.ParseUint(matched[1], 10, 16)
		viewportHeight, _ := strconv.ParseUint(matched[2], 10, 16)

		result.ViewportWidth = uint16(viewportWidth)
		result.ViewportHeight = uint16(viewportHeight)
		result.Viewport = fmt.Sprintf("%dx%d", result.ViewportWidth, result.ViewportHeight)

	}

	if sizeRegex.MatchString(screenSize) {
		matched := sizeRegex.FindStringSubmatch(screenSize)
		screenWidth, _ := strconv.ParseUint(matched[1], 10, 16)
		screenHeight, _ := strconv.ParseUint(matched[2], 10, 16)

		result.ScreenWidth = uint16(screenWidth)
		result.ScreenHeight = uint16(screenHeight)
		result.Screen = fmt.Sprintf("%dx%d", screenWidth, screenHeight)

		result.ResolutonWidth = uint16(float64(screenWidth) * pixelRatio)
		result.ResolutonHeight = uint16(float64(screenHeight) * pixelRatio)
		result.Resoluton = fmt.Sprintf("%dx%d", result.ResolutonWidth, result.ResolutonHeight)
	}

	return &result
}
