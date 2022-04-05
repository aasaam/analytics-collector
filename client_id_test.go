package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestClientIdFromStd1(t *testing.T) {
	initTime := time.Now().Add(time.Duration(-60) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	valid1 := base64.StdEncoding.EncodeToString([]byte(cid))
	fmt.Println(valid1)
	data1, err1 := clientIDStandardParser(valid1)
	if err1 != nil {
		t.Error(err1)
	}
	if data1.CidStdInitTime != time.Unix(initTime, 0) {
		t.Errorf("invalid init time")
	}
	if data1.CidStdSessionTime != time.Unix(sessionTime, 0) {
		t.Errorf("invalid session time")
	}
}

func TestClientIdFromStd2(t *testing.T) {
	initTime := time.Now().Add(time.Duration(-10) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
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
func TestClientIdFromStd7(t *testing.T) {
	it := time.Now().Add(time.Hour).Add(time.Second)
	st := time.Now()
	cid := strconv.Itoa(int(it.Unix())) + ":" + strconv.Itoa(int(st.Unix())) + ":0000000000000000"
	invalid1 := base64.StdEncoding.EncodeToString([]byte(cid))
	_, err1 := clientIDStandardParser(invalid1)
	if err1 == nil {
		t.Errorf("invalid base64 must be error")
	}
}
func TestClientIdFromNotStd1(t *testing.T) {
	c1 := clientIDFromAMP("amp-xyz")
	if len(c1.CidSessionChecksum) != 40 {
		t.Errorf("invalid cid")
	}
	c2 := clientIDFromOther([]string{"1"})
	if len(c2.CidSessionChecksum) != 40 {
		t.Errorf("invalid cid")
	}
}

func BenchmarkClientIdFromStd(b *testing.B) {
	initTime := time.Now().Add(time.Duration(-60) * time.Minute).Unix()
	sessionTime := time.Now().Add(time.Duration(-30) * time.Minute).Unix()
	cid := strconv.Itoa(int(initTime)) + ":" + strconv.Itoa(int(sessionTime)) + ":0000000000000000"
	valid1 := base64.StdEncoding.EncodeToString([]byte(cid))
	for n := 0; n < b.N; n++ {
		clientIDStandardParser(valid1)
	}
}
