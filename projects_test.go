package main

import (
	_ "embed"
	"encoding/json"
	"strconv"
	"testing"
)

//go:embed projects.json
var jsonSample []byte

func TestProjectManager0(t *testing.T) {
	if _, err := validatePublicInstaceID(""); err == nil {
		t.Errorf("error must throw")
	}
	if _, err := validatePublicInstaceID("000000000000"); err != nil {
		t.Errorf("error must not throw")
	}
}

func TestProjectManager1(t *testing.T) {
	pm := newProjectsManager()

	var data map[string]projectData

	err := json.Unmarshal(jsonSample, &data)
	if err != nil {
		t.Error(err)
	}

	err2 := pm.load(data)
	if err2 != nil {
		t.Error(err2)
	}

	if !pm.validateID("000000000000") {
		t.Errorf("project api must matched")
	}

	if pm.validateID("111111111111") {
		t.Errorf("project api must not matched")
	}

	if !pm.validateIDAndPrivate("000000000000", "000000000000111111111111") {
		t.Errorf("project api must matched")
	}

	if pm.validateIDAndPrivate("000000000000", "222222222222222222222222") {
		t.Errorf("project api must not matched")
	}

	if !pm.validateIDAndURL("000000000000", getURL("http://example.com")) {
		t.Errorf("project must matched")
	}

	if !pm.validateIDAndURL("000000000000", getURL("http://example.com")) {
		t.Errorf("project must matched")
	}

	if !pm.validateIDAndURL("000000000000", getURL("http://example.net")) {
		t.Errorf("project must matched")
	}

	if !pm.validateIDAndURL("000000000000", getURL("https://www.example.net")) {
		t.Errorf("project must matched")
	}

	if pm.validateIDAndURL("000000000000", getURL("http://1.1.1.1")) {
		t.Errorf("project must matched")
	}

	if pm.validateIDAndURL("000000000000", getURL("https://www.example-not-exist.net")) {
		t.Errorf("project must matched")
	}
}

func TestProjectManager2(t *testing.T) {
	pm := newProjectsManager()

	var data map[string]projectData

	err := json.Unmarshal(jsonSample, &data)
	if err != nil {
		t.Error(err)
	}

	err2 := pm.load(data)
	if err2 != nil {
		t.Error(err2)
	}

	if pm.validateIDAndURL("1", getURL("https://www.example.net")) {
		t.Errorf("project must matched")
	}
	if pm.validateIDAndURL("000000000000", getURL("")) {
		t.Errorf("project must matched")
	}
}

func BenchmarkValidateWildCardNoCache(b *testing.B) {
	pm := newProjectsManager()
	var data map[string]projectData
	json.Unmarshal(jsonSample, &data)
	pm.load(data)

	for n := 0; n < b.N; n++ {
		u := getURL("https://" + strconv.Itoa(n) + ".sub.example.net")
		pm.validateIDAndURL("000000000000", u)
	}
}

func BenchmarkValidateWildCardCache(b *testing.B) {
	pm := newProjectsManager()
	var data map[string]projectData
	json.Unmarshal(jsonSample, &data)
	pm.load(data)

	u := getURL("https://sub.example.net")
	for n := 0; n < b.N; n++ {
		pm.validateIDAndURL("000000000000", u)
	}
}

func BenchmarkValidateAPI(b *testing.B) {
	pm := newProjectsManager()
	var data map[string]projectData
	json.Unmarshal(jsonSample, &data)
	pm.load(data)

	for n := 0; n < b.N; n++ {
		pm.validateIDAndPrivate("000000000000", "000000000000111111111111")
	}
}
