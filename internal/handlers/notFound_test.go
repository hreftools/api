package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFount(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/404", nil)
	rr := httptest.NewRecorder()

	NotFound(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusNotFound)
	}

	var res NotFoundResponse

	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Errorf("error decoding response body: %v", err)
	}

	if res.Status != "error" {
		t.Errorf("unexpected status in response: got %v want %v", res.Status, "error")
	}

	if res.Data != "endpoint not found" {
		t.Errorf("unexpected data in response: got %v want %v", res.Data, "endpoint not found")
	}
}
