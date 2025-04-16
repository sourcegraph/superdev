package superdev

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendImageToServer(t *testing.T) {
	// Setup a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request is a POST to /run
		if r.Method != http.MethodPost || r.URL.Path != "/run" {
			t.Errorf("Expected POST /run request, got %s %s", r.Method, r.URL.Path)
		}

		// Check that the Content-Type is application/json
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		// Write a successful response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","message":"Request processed successfully","thread_id":"123456789"}`)) 
	}))
	defer server.Close()

	// Call the function to test
	threadID, err := sendImageToServer(server.URL, "superdev-wrapped-image", "test prompt")

	// Verify the results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if threadID != "123456789" {
		t.Errorf("Expected thread ID 123456789, got %s", threadID)
	}
}