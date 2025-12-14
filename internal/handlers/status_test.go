package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rr := httptest.NewRecorder()

	Status(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	var response StatusResponse

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("error decoding response body: %v", err)
	}
	if response.Status != "ok" {
		t.Errorf("unexpected status in response: got %v want %v", response.Status, "ok")
	}

	if response.Data != "service is running" {
		t.Errorf("unexpected data in response: got %v want %v", response.Data, "service is running")
	}

}
