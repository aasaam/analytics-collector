package main

import (
	"testing"
)

func TestRedis001(t *testing.T) {
	r, r1E := redisNew("redis://127.0.0.1/1")
	if r1E != nil {
		t.Error(r1E)
	}
	if r.Options().DB != 1 {
		t.Errorf("invalid db")
	}
	_, r2E := redisNew("redis://127.0.0.1")
	if r2E != nil {
		t.Error(r2E)
	}
}
