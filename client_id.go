package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	clientIDTypeOther uint8 = 0
	clientIDTypeStd   uint8 = 1
	clientIDTypeAmp   uint8 = 2
)

const minimumTime = 1640995200 // 2022-01-01 00:00:00

var cidStandardRegex = regexp.MustCompile(`^(?P<InitTime>[0-9]{9,11}):(?P<SessionTime>[0-9]{9,11}):(?P<Random>[a-z0-9]{16})$`)

type clientID struct {
	Valid              bool
	CidType            uint8
	CidUserChecksum    string
	CidSessionChecksum string
	CidStdInitTime     int64
	CidStdSessionTime  int64
}

func clientIDNoneSTD(parts []string, clientType uint8) clientID {
	if clientType == clientIDTypeStd {
		panic("none std must type none std")
	}

	t := time.Now()
	timestamp := t.Unix()
	sessionEach30Minutes := math.Floor(float64(timestamp) / 1800)
	sessionEach30MinutesString := fmt.Sprint(sessionEach30Minutes)

	userChecksum := checksum(strings.Join(parts, ":"))

	cid := clientID{
		Valid:              true,
		CidType:            clientType,
		CidUserChecksum:    userChecksum,
		CidSessionChecksum: checksum(userChecksum + ":" + sessionEach30MinutesString),
		CidStdInitTime:     0,
		CidStdSessionTime:  0,
	}

	return cid
}

func clientIDStandardParser(cidString string) (clientID, error) {
	cidInvalid := clientID{
		Valid: false,
	}

	decodedByte, err := base64.StdEncoding.DecodeString(cidString)
	if err != nil {
		return cidInvalid, err
	}

	decoded := string(decodedByte)
	if ok := cidStandardRegex.MatchString(decoded); ok {
		matched := cidStandardRegex.FindStringSubmatch(decoded)
		initTime, err := strconv.ParseInt(matched[1], 10, 64)
		if err != nil || initTime < minimumTime || initTime > time.Now().Add(time.Hour).Unix() {
			return cidInvalid, errors.New("invalid client identifier init time")
		}

		sessionTime, err := strconv.ParseInt(matched[2], 10, 64)
		if err != nil || sessionTime < initTime || sessionTime > time.Now().Add(time.Duration(12)*time.Hour).Unix() {
			return cidInvalid, errors.New("invalid client identifier session time")
		}

		diffSession := sessionTime - initTime
		if diffSession > 86400 {
			return cidInvalid, errors.New("session time must be at least 86400 with initialize time and past 24 hours")
		}

		cidValid := clientID{}
		cidValid.Valid = true
		cidValid.CidType = clientIDTypeStd
		cidValid.CidStdInitTime = time.Unix(initTime, 0).Unix()
		cidValid.CidStdSessionTime = time.Unix(sessionTime, 0).Unix()
		cidValid.CidUserChecksum = checksum(strings.Join([]string{matched[1], matched[3]}, ":"))
		cidValid.CidSessionChecksum = checksum(strings.Join([]string{matched[1], matched[2], matched[3]}, ":"))
		return cidValid, nil
	}

	return cidInvalid, errors.New("invalid client identifier")
}
