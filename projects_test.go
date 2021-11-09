package main

import (
	"strconv"
	"testing"
)

var jsonSample = `
{
	"000000000000": {
		"ph": "000000000000",
		"d": [
			"example.com",
			"*.sub.example.com"
		]
	},
	"111111111111": {
		"ph": "111111111111",
		"d": [
			"example.net",
			"*.sub.example.net"
		]
	}
}
`

func TestProjectManager1(t *testing.T) {
	pManager := NewProjectManager()

	err1 := pManager.LoadJSON([]byte(`1`))
	if err1 == nil {
		t.Errorf("error must be thrown")
	}

	err := pManager.LoadJSON([]byte(jsonSample))
	if err != nil {
		t.Error(err)
	}

	v00 := pManager.ValidateEvent("000000000000")
	if !v00 {
		t.Errorf("project api must matched")
	}

	v0 := pManager.ValidateAPI("000000000000", "000000000000")
	if !v0 {
		t.Errorf("project api must matched")
	}

	v0f := pManager.ValidateAPI("000000000000", "000000000001")
	if v0f {
		t.Errorf("project api must not match")
	}

	v1 := pManager.ValidatePageView("000000000001", "example.com")
	if v1 {
		t.Errorf("project id must not matched")
	}

	v31 := pManager.ValidatePageView("000000000000", "example.com")
	v32 := pManager.ValidatePageView("000000000000", "example.com")
	if !v31 || !v32 {
		t.Errorf("project id must matched")
	}

	v2 := pManager.ValidatePageView("0", "yahoo.com")
	if v2 {
		t.Errorf("project id must not matched")
	}

	v41 := pManager.ValidatePageView("000000000000", "very.sub.example.com")
	v42 := pManager.ValidatePageView("000000000000", "sub.example.com")
	if !v41 || !v42 {
		t.Errorf("project id must matched")
	}
}

func BenchmarkValidateWildCardNoCache(b *testing.B) {
	pManager := NewProjectManager()
	pManager.LoadJSON([]byte(jsonSample))
	for n := 0; n < b.N; n++ {
		pManager.ValidatePageView("000000000000", strconv.Itoa(n)+".sub.example.com")
	}
}
func BenchmarkValidateWildCardCache(b *testing.B) {
	pManager := NewProjectManager()
	pManager.LoadJSON([]byte(jsonSample))
	for n := 0; n < b.N; n++ {
		pManager.ValidatePageView("000000000000", "static.sub.example.com")
	}
}
func BenchmarkValidateAPI(b *testing.B) {
	pManager := NewProjectManager()
	pManager.LoadJSON([]byte(jsonSample))
	for n := 0; n < b.N; n++ {
		pManager.ValidateAPI("000000000000", "000000000001")
	}
}
