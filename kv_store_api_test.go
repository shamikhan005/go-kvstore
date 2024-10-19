package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlePost(t *testing.T) {
	reqBody := strings.NewReader(`{"key": "winter", "value": "coming"}`)
	req, _ := http.NewRequest("POST", "/kvstore", reqBody)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(keyValueStoreHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestHandleGet(t *testing.T) {
	kvs.Put("winter", "coming")

	req, _ := http.NewRequest("GET", "/kvstore?key=winter", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(keyValueStoreHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}