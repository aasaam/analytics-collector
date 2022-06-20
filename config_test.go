package main

import (
	"testing"
)

func TestConfig(t *testing.T) {
	c1 := newConfig("trace", 0, true, "http://127.0.0.1")
	c1.getLogger().Warn().Msg("Warn")
	c1.getLogger().Trace().Msg("Trace")

	if c1.getLogger() == nil {
		t.Errorf("no logger")
	}

	c2 := newConfig("error", 0, true, "http://127.0.0.1")
	c2.getLogger().Warn().Msg("Warn")
	c2.getLogger().Trace().Msg("Trace")

}
