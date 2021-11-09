package main

import (
	"testing"
)

func TestClientIdFromStd(t *testing.T) {
	if !IsAmpClient("amp-w1YG5GrNZrbO6LsZ4EzWtA") {
		t.Errorf("amp not detect")
	}
	valid1 := "YzoxNTc3ODI0MjAwOjE1Nzc4MjQyMDA6MWZqdjhhdDhva2xxaG8wcw=="
	data1, err1 := ClientIdFromStd(valid1)
	if err1 != nil || data1.StdType != "c" {
		t.Error(err1)
	}
	prettyPrint(data1)

	valid2 := "dDoxNTc3ODI0MjAwOjE1Nzc4MjQyMDA6MWZqdjhhdDhva2xxaG8wcw=="
	data2, err2 := ClientIdFromStd(valid2)
	if err2 != nil || data2.StdType != "t" {
		t.Error(err1)
	}
	prettyPrint(data2)

	invalid1 := "YzowMDAwMDAwMDAwOjAwMDAwMDAwMDA6MWZqdjhhdDhva2xxaG8wcw=="
	_, ier2 := ClientIdFromStd(invalid1)
	if ier2 == nil {
		t.Errorf("error must be thrown")
	}
	invalid2 := "YzoxNTc3ODI0MjAwOjE1Nzc4MjQxMDA6MWZqdjhhdDhva2xxaG8wcw=="
	_, ier3 := ClientIdFromStd(invalid2)
	if ier3 == nil {
		t.Errorf("error must be thrown")
	}
}
func TestClientIdFromHashParts(t *testing.T) {
	hashParts := []string{"127.0.0.1", "curl/7"}
	cid1 := ClientIdFromHashParts(ClientTypeAmp, hashParts)
	if cid1.Hash == "" {
		t.Errorf("invalid cliend id")
	}
	cid2 := ClientIdFromHashParts(ClientTypeAmp, hashParts)
	if cid1.Hash != cid2.Hash {
		t.Errorf("invalid cliend id same params")
	}

	hashParts = []string{"127.0.0.2", "curl/7"}
	ocid1 := ClientIdFromHashParts(ClientTypeAmp, hashParts)
	prettyPrint(ocid1)
	if cid1.Hash == ocid1.Hash {
		t.Errorf("invalid cliend id same params")
	}
}
