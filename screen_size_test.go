package main

import (
	"testing"
)

func TestParseScreenSize(t *testing.T) {
	scr1 := ParseScreenSize("p-p", "1376x774", "1376x389", 1.3953488372093024, 24)

	if scr1.Resoluton != "1920x1080" || scr1.ScreenOrientation != ScreenOrientationPortraitPrimary {
		t.Errorf("invalid screen orientation")
	}
	scr11 := ParseScreenSize("p-s", "1376x774", "1376x389", 1.3953488372093024, 24)

	if scr11.ScreenOrientation != ScreenOrientationPortraitSecondary {
		t.Errorf("invalid screen orientation")
	}

	scr12 := ParseScreenSize("l-p", "1376x774", "1376x389", 1.3953488372093024, 24)

	if scr12.ScreenOrientation != ScreenOrientationLandscapePrimary {
		t.Errorf("invalid screen orientation")
	}

	scr13 := ParseScreenSize("l-s", "1376x774", "1376x389", 1.3953488372093024, 24)
	if scr13.ScreenOrientation != ScreenOrientationLandscapeSecondary {
		t.Errorf("invalid screen orientation")
	}

	if scr12.Resoluton != "1920x1080" {
		t.Errorf("invalid resolution")
	}

	scr2 := ParseScreenSize("l-p", "2x2", "1x1", 1.3953488372093024, 24)

	if scr2.ColorDepth != 24 {
		t.Errorf("invalid resolution")
	}

}

func BenchmarkParseScreenSize(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ParseScreenSize("l-p", "1376x774", "1376x389", 1.3953488372093024, 24)
	}
}
