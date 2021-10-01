package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	recorder := httptest.NewRecorder()
	handleHealthcheck(recorder, request)
	if recorder.Code != 200 || recorder.Body.String() != "Okay" {
		t.Fatalf("Unexpected %v %v", recorder.Code, recorder.Body.String())
	}
}
