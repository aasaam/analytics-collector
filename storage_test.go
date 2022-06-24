package main

import (
	"testing"
)

func TestStorage1(t *testing.T) {
	storage := newStorage()

	storage.addRecord([]byte("1"))
	storage.addRecord([]byte("2"))
	storage.addRecord([]byte("3"))

	if storage.recordCount != 3 {
		t.Errorf("invalid length")
	}

	if storage.recordCount != 3 {
		t.Errorf("invalid length")
	}

	storage.cleanRecords()

	if storage.recordCount != 0 {
		t.Errorf("invalid length")
	}

	storage.addClientError([]byte("1"))
	storage.addClientError([]byte("2"))

	if storage.clientErrorCount != 2 {
		t.Errorf("invalid length")
	}

	if storage.clientErrorCount != 2 {
		t.Errorf("invalid length")
	}

	storage.cleanClientErrors()

	if storage.clientErrorCount != 0 {
		t.Errorf("invalid length")
	}
}

func BenchmarkStorage1(b *testing.B) {
	storage := newStorage()

	for n := 0; n < b.N; n++ {
		storage.addRecord([]byte("1"))
	}
}
