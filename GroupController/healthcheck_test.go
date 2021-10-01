package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	controller, _ := NewController(nil, context.Background())
	request := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	recorder := httptest.NewRecorder()
	controller.handleHealthcheck(recorder, request)
	if recorder.Code != 200 || recorder.Body.String() != "Okay" {
		t.Fatalf("Unexpected %v %v", recorder.Code, recorder.Body.String())
	}
}
