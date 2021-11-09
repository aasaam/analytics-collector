package main

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"regexp"
	"strconv"
	"time"
)

const (
	// ClientTypeOther other types
	ClientTypeOther uint8 = 0
	// ClientTypeStd standard js client algorithm
	ClientTypeStd uint8 = 1
	// ClientTypeAmp for amp client
	ClientTypeAmp uint8 = 2
	// ClientTypeImg for none script and image request
	ClientTypeImg uint8 = 3
)

const minimumTime = 1577824200 // 2020-01-01 00:00:00

var cidRegex, _ = regexp.Compile(`^(?P<TypeOfClientID>[a-z]{1}):(?P<InitTime>[0-9]{9,11}):(?P<DailyTime>[0-9]{9,11}):(?P<Random>[a-z0-9]{16})$`)

var ampRegex, _ = regexp.Compile("^amp-(.*)$")

var hashFixRegex, _ = regexp.Compile("[^a-zA-Z0-9]+")

// ClientIdentifier is client hash for tracking
type ClientIdentifier struct {
	Type uint8

	Hash string

	// standard client hash
	StdType      string
	StdInitTime  time.Time
	StdRandom    string
	StdDailyTime time.Time
}

// HashFromPartsForClientId will return base 62 fix 16 bytes
func HashFromPartsForClientId(parts []string) string {
	h := sha1.New()
	for _, s := range parts {
		io.WriteString(h, s)
	}
	return hashFixRegex.ReplaceAllString(base64.StdEncoding.EncodeToString(h.Sum(nil)), "")[0:16]
}

// IsAmpClient is matched for amp client
func IsAmpClient(cid string) bool {
	return ampRegex.MatchString(cid)
}

// ClientIdFromHashParts will store hash of client static parts
func ClientIdFromHashParts(typeOfCid uint8, parts []string) *ClientIdentifier {

	cid := ClientIdentifier{
		Type: typeOfCid,
		Hash: HashFromPartsForClientId(parts),
	}

	return &cid
}

// ClientIdFromStd is standard JavaScript client identifier parser
func ClientIdFromStd(cidString string) (*ClientIdentifier, error) {
	cid := ClientIdentifier{}
	decodedByte, err := base64.StdEncoding.DecodeString(cidString)
	if err != nil {
		return &cid, err
	}
	decoded := string(decodedByte)
	if ok := cidRegex.MatchString(decoded); ok {
		matched := cidRegex.FindStringSubmatch(decoded)
		InitTime, err := strconv.ParseInt(matched[2], 10, 64)
		if err != nil || InitTime < minimumTime {
			return &cid, errors.New("invalid client identifier init time")
		}
		DailyTime, err := strconv.ParseInt(matched[3], 10, 64)
		if err != nil || DailyTime < InitTime {
			return &cid, errors.New("invalid client identifier daily time")
		}
		cid.Type = ClientTypeStd
		cid.StdType = matched[1]
		cid.StdInitTime = time.Unix(InitTime, 0)
		cid.StdDailyTime = time.Unix(DailyTime, 0)
		cid.StdRandom = matched[3]
		cid.Hash = HashFromPartsForClientId([]string{matched[1], matched[2], matched[4]})
		return &cid, nil
	}
	return &cid, errors.New("invalid client identifier")
}
