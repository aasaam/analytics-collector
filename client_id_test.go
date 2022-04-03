package main

import (
	"encoding/base64"
	"testing"
	"time"
)

func TestClientIdFromStd1(t *testing.T) {
	cid := "1647253200:1647253300:0000000000000000"
	valid1 := base64.StdEncoding.EncodeToString([]byte(cid))
	data1, err1 := clientIDStandardParser(valid1)
	if err1 != nil {
		t.Error(err1)
	}
	if data1.CIDStdInitTime != time.Unix(1647253200, 0) {
		t.Errorf("invalid init time")
	}
	if data1.CIDStdSessionTime != time.Unix(1647253300, 0) {
		t.Errorf("invalid session time")
	}
}

func TestClientIdFromStd2(t *testing.T) {
	cid := "1647253300:1647253200:0000000000000000"
	invalid1 := base64.StdEncoding.EncodeToString([]byte(cid))
	_, err1 := clientIDStandardParser(invalid1)
	if err1 == nil {
		t.Errorf("init time must be before session time")
	}
}
func TestClientIdFromStd3(t *testing.T) {
	cid := "1640995140:1647253300:0000000000000000"
	invalid1 := base64.StdEncoding.EncodeToString([]byte(cid))
	_, err1 := clientIDStandardParser(invalid1)
	if err1 == nil {
		t.Errorf("init time must start with 2022")
	}
}

func TestClientIdFromStd4(t *testing.T) {
	cid := "1647253200:1647339601:0000000000000000"
	invalid1 := base64.StdEncoding.EncodeToString([]byte(cid))
	_, err1 := clientIDStandardParser(invalid1)
	if err1 == nil {
		t.Errorf("init time must start with 2022")
	}
}
func TestClientIdFromStd5(t *testing.T) {
	cid := "1647253200-1647339601-0000000000000000"
	invalid1 := base64.StdEncoding.EncodeToString([]byte(cid))
	_, err1 := clientIDStandardParser(invalid1)
	if err1 == nil {
		t.Errorf("invalid pattern must be error")
	}
}
func TestClientIdFromStd6(t *testing.T) {
	invalid1 := "./-()*"
	_, err1 := clientIDStandardParser(invalid1)
	if err1 == nil {
		t.Errorf("invalid base64 must be error")
	}
}
func TestClientIdFromNotStd1(t *testing.T) {
	c1 := clientIDFromAMP("amp-xyz")
	if len(c1.CIDSessionChecksum) != 40 {
		t.Errorf("invalid cid")
	}
	c2 := clientIDFromOther([]string{"1"})
	if len(c2.CIDSessionChecksum) != 40 {
		t.Errorf("invalid cid")
	}
}

func BenchmarkClientIdFromStd(b *testing.B) {
	cid := "1647253200:1647253300:0000000000000000"
	valid1 := base64.StdEncoding.EncodeToString([]byte(cid))
	for n := 0; n < b.N; n++ {
		clientIDStandardParser(valid1)
	}
}
