package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllocate(t *testing.T) {
	allocateRequest := AllocateRequest{
		Id:     "Suriya",
		Target: 255,
	}
	body, _ := json.Marshal(allocateRequest)
	request := httptest.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(body))
	recorder := httptest.NewRecorder()
	handleAllocate(recorder, request)
	if recorder.Body.String() != "Suriya acked with the target of ff" {
		t.Fatalf("Unexpected response %v", recorder.Body.String())
	}
}
