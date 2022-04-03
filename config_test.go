package main

import (
	"net"
	"testing"
)

func TestConfig(t *testing.T) {
	ip := net.ParseIP("1.1.1.1")
	c1 := newConfig("trace", 0, true, "http://127.0.0.1", "1.1.1.1,8.8.8.8")
	c1.getLogger().Warn().Msg("Warn")
	c1.getLogger().Trace().Msg("Trace")

	if !c1.canAccessMetrics(ip) {
		t.Errorf("ip validation failed")
	}

	if c1.getLogger() == nil {
		t.Errorf("no logger")
	}

	c2 := newConfig("error", 0, true, "http://127.0.0.1", "")
	c2.getLogger().Warn().Msg("Warn")
	c2.getLogger().Trace().Msg("Trace")

	if c2.canAccessMetrics(ip) {
		t.Errorf("ip validation failed")
	}
}
