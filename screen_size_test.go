package main

import (
	"testing"
)

func TestParseScreenSize(t *testing.T) {
	scr1 := parseScreenSize("portrait-primary", "1376x774", "1376x389", "1.3953488372093024", "24")

	if scr1.ScrColorDepth != 24 {
		t.Errorf("invalid color depth")
	}

	if scr1.ScrDevicePixelRatio != 1.3953488372093024 {
		t.Errorf("invalid device pixel ratio")
	}

	if scr1.ScrResolution != "1920x1080" || !scr1.ScrScreenOrientationIsPortrait {
		t.Errorf("invalid screen orientation")
	}

	if scr1.ScrResolution == "" {
		t.Errorf("invalid Resolution")
	}

	scr11 := parseScreenSize("portrait-secondary", "1376x774", "1376x389", "1", "16")

	if scr11.ScrColorDepth != 16 {
		t.Errorf("invalid color depth")
	}

	if scr11.ScrDevicePixelRatio != 1 {
		t.Errorf("invalid device pixel ratio")
	}

	if !scr11.ScrScreenOrientationIsPortrait || !scr11.ScrScreenOrientationIsSecondary {
		t.Errorf("invalid screen orientation")
	}

	scr12 := parseScreenSize("landscape-primary", "1376x774", "1376x389", "", "24")

	if scr12.ScrResolution != "" {
		t.Errorf("invalid resolution")
	}

	if scr12.ScrScreenOrientationIsPortrait || scr12.ScrScreenOrientationIsSecondary {
		t.Errorf("invalid screen orientation")
	}

}

func TestParseScreenSize2(t *testing.T) {
	scr1 := parseScreenSize("landscape-secondary", "1376x774", "1376x389", "1.3953488372093024", "24")

	if scr1.ScrColorDepth != 24 {
		t.Errorf("invalid color depth")
	}

	if scr1.ScrDevicePixelRatio != 1.3953488372093024 {
		t.Errorf("invalid device pixel ratio")
	}

	if scr1.ScrResolution != "1920x1080" || scr1.ScrScreenOrientationIsPortrait || !scr1.ScrScreenOrientationIsSecondary {
		t.Errorf("invalid screen orientation")
	}

	if scr1.ScrResolution == "" {
		t.Errorf("invalid resolution")
	}

}

func BenchmarkParseScreenSize(b *testing.B) {
	for n := 0; n < b.N; n++ {
		parseScreenSize("landscape-primary", "1376x774", "1376x389", "1.3953488372093024", "24")
	}
}
