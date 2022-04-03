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
	}

	prometheusStorageRecords.Add(0)

	return &s
}

func (s *storage) setRecords(records [][]byte) {
	s.records = records
	s.recordCount = len(records)
	prometheusStorageRecords.Add(float64(s.recordCount))
}

func (s *storage) getRecords() [][]byte {
	s.Lock()
	defer s.Unlock()
	return s.records
}

func (s *storage) cleanRecords() {
	records := make([][]byte, 0)
	s.records = records
	s.recordCount = 0
	prometheusStorageRecords.Set(float64(s.recordCount))
}

func (s *storage) addRecord(r []byte) {
	s.Lock()
	defer s.Unlock()
	s.records = append(s.records, r)
	s.recordCount++
	prometheusStorageRecords.Set(float64(s.recordCount))
}

func (s *storage) addClientError(r []byte) {
	s.Lock()
	defer s.Unlock()
	s.records = append(s.records, r)
	s.recordCount++
	prometheusStorageRecords.Set(float64(s.recordCount))
}

// package main

// import (
// 	"sync"
// )

// type storage struct {
// 	sync.Mutex
// 	count int
// 	items []record
// }

// func newStorage() *storage {
// 	s := storage{
// 		count: 0,
// 	}

// 	prometheusStorageItems.Add(0)

// 	return &s
// }

// func (s *storage) setItems(items []record) {
// 	s.items = items
// 	s.count = len(items)
// 	prometheusStorageItems.Add(float64(s.count))
// }

// func (s *storage) getItems() []record {
// 	s.Lock()
// 	defer s.Unlock()
// 	return s.items
// }

// func (s *storage) clean() {
// 	items := make([]record, 0)
// 	s.items = items
// 	s.count = 0
// 	prometheusStorageItems.Set(float64(s.count))
// }

// func (s *storage) addRecord(r record) {
// 	s.Lock()
// 	defer s.Unlock()
// 	s.items = append(s.items, r)
// 	s.count++
// 	prometheusStorageItems.Set(float64(s.count))
// }
