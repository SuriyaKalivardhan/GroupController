package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllocate(t *testing.T) {
	controller, _ := NewController(nil, context.Background())
	allocateRequest := AllocateRequest{
		Id:     "Suriya",
		Target: 255,
	}
	body, _ := json.Marshal(allocateRequest)
	request := httptest.NewRequest(http.MethodPost, "/allocate", bytes.NewReader(body))
	recorder := httptest.NewRecorder()
	controller.handleAllocate(recorder, request)
	if recorder.Body.String() != "Suriya acked with the target of ff" {
		t.Fatalf("Unexpected response %v", recorder.Body.String())
	}
}
