package main

import "testing"

func TestKeyValueStore(t *testing.T) {
	kvs := NewKeyValueStore()

	kvs.Put("testKey", "testValue")
	value, exists := kvs.Get("testKey")
	if !exists || value != "testValue" {
		t.Errorf("Expected 'testValue', got '%s'", value)
	}

	kvs.Delete("testKey")
	_, exists = kvs.Get("testKey")
	if exists {
		t.Error("Expected key to be deleted")
	}
}
