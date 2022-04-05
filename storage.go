package main

import (
	"sync"
)

type storage struct {
	sync.Mutex
	clientErrorCount int
	clientErrors     [][]byte
	recordCount      int
	records          [][]byte
}

func newStorage() *storage {
	s := storage{
		clientErrorCount: 0,
		recordCount:      0,
		clientErrors:     make([][]byte, 0),
		records:          make([][]byte, 0),
	}

	promMetricStorageQueueRecords.Set(0)
	promMetricStorageQueueClientErrors.Set(0)

	return &s
}

func (s *storage) setRecords(items [][]byte) {
	s.records = items
	s.recordCount = len(items)
	promMetricStorageQueueRecords.Set(float64(s.recordCount))
}

func (s *storage) getRecords() [][]byte {
	return s.records
}

func (s *storage) addRecord(r []byte) {
	s.Lock()
	defer s.Unlock()
	s.records = append(s.records, r)
	s.recordCount++
	promMetricStorageQueueRecords.Set(float64(s.recordCount))
}

func (s *storage) cleanRecords() {
	items := make([][]byte, 0)
	s.records = items
	s.recordCount = 0
	promMetricStorageQueueRecords.Set(0)
}

func (s *storage) setClientErrors(items [][]byte) {
	s.clientErrors = items
	s.clientErrorCount = len(items)
	promMetricStorageQueueClientErrors.Set(float64(s.recordCount))
}

func (s *storage) getClientErrors() [][]byte {
	return s.clientErrors
}

func (s *storage) cleanClientErrors() {
	items := make([][]byte, 0)
	s.clientErrors = items
	s.clientErrorCount = 0
	promMetricStorageQueueClientErrors.Set(0)
}

func (s *storage) addClientError(r []byte) {
	s.Lock()
	defer s.Unlock()
	s.clientErrors = append(s.clientErrors, r)
	s.clientErrorCount++
	promMetricStorageQueueRecords.Set(float64(s.clientErrorCount))
}
