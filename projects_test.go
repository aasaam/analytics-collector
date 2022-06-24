package main

import (
	_ "embed"
	"encoding/json"
	"strconv"
	"testing"
)

//go:embed projects.json
var jsonSample []byte

func getTestProjects() *projects {
	projectsManager := newProjectsManager()

	var data map[string]projectData

	err := json.Unmarshal(jsonSample, &data)
	if err != nil {
		panic(err)
	}

	err2 := projectsManager.load(data)
	if err2 != nil {
		panic(err2)
	}

	return projectsManager
}

func TestProjectManagerJSONFile(t *testing.T) {
	projectsLoad("./projects.json")
}

func TestProjectManager0(t *testing.T) {
	if _, err := validatePublicInstanceID(""); err == nil {
		t.Errorf("error must throw")
	}
	if _, err := validatePublicInstanceID("000000000000"); err != nil {
		t.Errorf("error must not throw")
	}
}

func TestProjectManager1(t *testing.T) {
	pm := getTestProjects()

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
	pm := getTestProjects()

	if pm.validateIDAndURL("1", getURL("https://www.example.net")) {
		t.Errorf("project must matched")
	}
	if pm.validateIDAndURL("000000000000", getURL("")) {
		t.Errorf("project must matched")
	}
}
func TestProjectManager3(t *testing.T) {
	_, psErr := projectsLoadJSON("./projects.json")

	if psErr != nil {
		t.Error(psErr)
	}
}

func BenchmarkValidateWildCardNoCache(b *testing.B) {
	pm := getTestProjects()

	for n := 0; n < b.N; n++ {
		u := getURL("https://" + strconv.Itoa(n) + ".sub.example.net")
		pm.validateIDAndURL("000000000000", u)
	}
}

func BenchmarkValidateWildCardCache(b *testing.B) {
	pm := getTestProjects()

	u := getURL("https://sub.example.net")
	for n := 0; n < b.N; n++ {
		pm.validateIDAndURL("000000000000", u)
	}
}

func BenchmarkValidateAPI(b *testing.B) {
	pm := getTestProjects()

	for n := 0; n < b.N; n++ {
		pm.validateIDAndPrivate("000000000000", "000000000000111111111111")
	}
}
