package main

import (
	"testing"
)

func TestRecord0(t *testing.T) {
	_, e1 := redisGetClient("redis://127.0.0.1:80/0")
	if e1 == nil {
		t.Errorf("must be invalid")
	}
	_, e2 := redisGetClient("!")
	if e2 == nil {
		t.Errorf("must be invalid")
	}
}
