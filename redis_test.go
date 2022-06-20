package main

import (
	"testing"
)

func TestRedis001(t *testing.T) {
	r, r1E := redisClientNew("redis://127.0.0.1:6379/1", 1)
	if r1E != nil {
		t.Error(r1E)
	}
	if r.rdb.Options().DB != 1 {
		t.Errorf("invalid db")
	}
	_, r2E := redisClientNew("redis://127.0.0.1:6379", 1)
	if r2E != nil {
		t.Error(r2E)
	}
}
func TestRedis002(t *testing.T) {
	r, r1E := redisClientNew("redis://127.0.0.1:6379/1", 5)
	if r1E != nil {
		t.Error(r1E)
	}
	v := "v1"
	v1 := []byte(v)
	r.pushRecord(v1)
	r.pushClientError(v1)

	i, ie := r.popRecord()
	r.popRecordSubmit()

	if ie != nil {
		t.Error(ie)
	}

	if len(i) != 1 || i[0] != v {
		t.Errorf("invalid response")
	}

}

func BenchmarkRedis1(b *testing.B) {
	r, _ := redisClientNew("redis://127.0.0.1:6379/1", 5)
	a := []byte("b")
	for n := 0; n < b.N; n++ {
		r.pushRecord(a)
	}
}
func BenchmarkRedis2(b *testing.B) {
	r, _ := redisClientNew("redis://127.0.0.1:6379/1", 5)
	for n := 0; n < b.N; n++ {
		r.popRecord()
		r.popRecordSubmit()
	}
}
