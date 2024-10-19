package main 

import (
	"testing"
)

func BenchmarkPut(b *testing.B) {
	kvs := NewKeyValueStore()
	for i := 0; i < b.N; i++ {
		kvs.Put("key", "value")
	}
}

func BenchmarkGet(b *testing.B) {
	kvs := NewKeyValueStore()
	kvs.Put("key", "value")
	for i := 0; i < b.N; i++ {
		kvs.Get("key")
	}
}

func BenchmarkDelete(b *testing.B) {
	kvs := NewKeyValueStore()
	kvs.Put("key", "value") 
	for i := 0; i < b.N; i++ {
		kvs.Delete("key")
	}
}