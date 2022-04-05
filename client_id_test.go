package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestClientIdFromStd0(t *testing.T) {
	d := `{"cid_std":"MTY0OTEwOTcwMToxNjQ5MTExNTAxOjAwMDAwMDAwMDAwMDAwMDA=","p":{"u":"https://exmple.net/","t":"نرم افزارها جامعه کتابهای ایپسوم با، نامفهوم، می طلبد و زمان گرافیک در","l":"fa","cu":"https://exmple.net/","ei":"63","em":"page","et":"A0000","r":"https://www.google.com","bc":{"n1":"ارائه آنچنان-1","u1":"/1-%D8%A7%DB%8C%D8%AC%D8%A7%D8%AF-%D9%85%DB%8C-%D8%B7%D9%84%D8%A8%D8%AF-%D9%88-%DA%A9%D8%A7%D8%B1%D8%A8%D8%B1%D8%AF%D9%87%D8%A7%DB%8C-%DA%86%D8%A7%D9%BE-%DA%A9%D8%AA%D8%A7%D8%A8%D9%87%D8%A7%DB%8C","n2":"و دشواری که-2","u2":"/2-%D9%88-%D8%AF%D8%B4%D9%88%D8%A7%D8%B1%DB%8C-%D9%85%D8%AA%D9%86-%D8%A7%D8%B3%D8%AA%D9%81%D8%A7%D8%AF%D9%87-%D9%BE%DB%8C%D9%88%D8%B3%D8%AA%D9%87-%D8%B3%D8%AA%D9%88%D9%86"},"scr":"393x851","vps":"405x740","cd":"24","k":"مورد,نیاز,امید,مورد,برای,فراوان,قرار,گیرد,شرایط,روزنامه","rs":"","dpr":"2.75","if":false,"ts":true,"sot":"p-p","prf":{"dlt":"0","tct":"0","srt":"111","pdt":"9","rt":"27","dit":"383","clt":"383","r":7},"geo":{"lat":35.722035,"lon":51.4074033}}}`
	var s postRequest
	e := json.Unmarshal([]byte(d), &s)
	fmt.Println(e)
}
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
