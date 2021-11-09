package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

func prettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}

func TestParseUTM(t *testing.T) {
	sampleName := " This IS <>\"' for Name"

	utm0 := ParseUTM("\x18")

	if utm0.Valid {
		t.Errorf("invalid utm parse")
	}

	utm1 := ParseUTM("https://www.google.com")

	if utm1.Valid {
		t.Errorf("invalid utm parse")
	}

	utm11 := ParseUTM("https://www.example.com?UTM_SOURCE=source&utm_medium=medium")

	if utm11.Valid {
		t.Errorf("invalid utm parse")
	}

	utm2 := ParseUTM("https://www.example.com?UTM_SOURCE=source&utm_medium=medium&utm_campaign=sale1&utm_id=id&utm_term=keyword1&utm_content=content")

	if !utm2.Valid || utm2.CampaignName == "" {
		t.Errorf("invalid utm parse")
	}

	utm3 := ParseUTM("https://www.example.com?UTM_SOURCE=source&utm_medium=medium&utm_campaign=" + sampleName + "&utm_id=id&utm_term=keyword1&utm_content=content")

	if !utm3.Valid {
		t.Errorf("invalid utm parse")
	}

	utm4 := ParseUTM("https://www.example.com/?UTM_SOURCE=source&utm_medium=medium&utm_campaign=" + url.QueryEscape(sampleName) + "&utm_id=id&utm_term=keyword1&utm_content=content")

	if !utm4.Valid {
		t.Errorf("invalid utm parse")
	}

	if utm4.CampaignName != utm3.CampaignName {
		t.Errorf("invalid utm parse")
	}
}
