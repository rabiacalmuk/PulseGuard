package checks

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPCheck_HealthyService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := HTTPCheck{URL: server.URL, ExpectedStatus: http.StatusOK}
	result, err := c.Run()
	if err != nil {
		t.Fatalf("Run hata verdi: %v", err)
	}
	if result.Level != LevelInfo {
		t.Errorf("Level yanlis: %s, beklenen INFO", result.Level)
	}
}

func TestHTTPCheck_UnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := HTTPCheck{URL: server.URL, ExpectedStatus: http.StatusOK}
	result, err := c.Run()
	if err != nil {
		t.Fatalf("Run hata verdi: %v", err)
	}
	if result.Level != LevelError {
		t.Errorf("Level yanlis: %s, beklenen ERROR", result.Level)
	}
}

func TestHTTPCheck_Unreachable(t *testing.T) {
	c := HTTPCheck{URL: "http://localhost:1", ExpectedStatus: http.StatusOK}
	result, err := c.Run()
	if err != nil {
		t.Fatalf("Run hata verdi: %v", err)
	}
	if result.Level != LevelError {
		t.Errorf("Level yanlis: %s, beklenen ERROR", result.Level)
	}
}
