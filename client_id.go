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
	CIDType            uint8
	CIDUserChecksum    string
	CIDSessionChecksum string
	CIDStdInitTime     time.Time
	CIDStdSessionTime  time.Time
}

func clientIDFromAMP(ampCIDString string) clientID {
	t := time.Now()
	timestamp := t.Unix()
	sessionEach30Minutes := math.Round(float64(timestamp) / 1800)
	sessionEach30MinutesString := fmt.Sprint(sessionEach30Minutes)

	userChecksum := checksum(ampCIDString)

	cid := clientID{
		Valid:              true,
		CIDType:            clientIDTypeAmp,
		CIDUserChecksum:    userChecksum,
		CIDSessionChecksum: checksum(userChecksum + ":" + sessionEach30MinutesString),
	}

	return cid
}

func clientIDFromOther(parts []string) clientID {
	t := time.Now()
	timestamp := t.Unix()
	sessionEach30Minutes := math.Round(float64(timestamp) / 1800)
	sessionEach30MinutesString := fmt.Sprint(sessionEach30Minutes)

	userChecksum := checksum(strings.Join(parts, ":"))

	cid := clientID{
		Valid:              true,
		CIDType:            clientIDTypeOther,
		CIDUserChecksum:    userChecksum,
		CIDSessionChecksum: checksum(userChecksum + ":" + sessionEach30MinutesString),
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
		if err != nil || initTime < minimumTime {
			return cidInvalid, errors.New("invalid client identifier init time")
		}
		sessionTime, err := strconv.ParseInt(matched[2], 10, 64)
		if err != nil || sessionTime < initTime {
			return cidInvalid, errors.New("invalid client identifier daily time")
		}
		if (sessionTime - initTime) > 86400 {
			return cidInvalid, errors.New("session time must be at least 86400 with initialize time")
		}
		cidValid := clientID{}
		cidValid.Valid = true
		cidValid.CIDType = clientIDTypeStd
		cidValid.CIDStdInitTime = time.Unix(initTime, 0)
		cidValid.CIDStdSessionTime = time.Unix(sessionTime, 0)
		cidValid.CIDUserChecksum = checksum(strings.Join([]string{matched[1], matched[3]}, ":"))
		cidValid.CIDSessionChecksum = checksum(strings.Join([]string{matched[1], matched[2], matched[3]}, ":"))
		return cidValid, nil
	}
	return cidInvalid, errors.New("invalid client identifier")
}
