package checks

import (
	"fmt"
	"net/http"
	"time"
)

type HTTPCheck struct {
	URL            string
	ExpectedStatus int
}

func (h HTTPCheck) Name() string {
	return "http_check"
}

func (h HTTPCheck) Run() (CheckResult, error) {
	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(h.URL)
	if err != nil {
		return CheckResult{
			CheckType: "http_check",
			Level:     LevelError,
			Message:   fmt.Sprintf("%s adresine ulasilamadi: %v", h.URL, err),
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != h.ExpectedStatus {
		return CheckResult{
			CheckType: "http_check",
			Level:     LevelError,
			Message:   fmt.Sprintf("%s beklenmeyen durum kodu dondu: %d (beklenen: %d)", h.URL, resp.StatusCode, h.ExpectedStatus),
		}, nil
	}

	return CheckResult{
		CheckType: "http_check",
		Level:     LevelInfo,
		Message:   fmt.Sprintf("%s saglikli (HTTP %d)", h.URL, resp.StatusCode),
	}, nil
}
