package api

import (
	"encoding/json"
	"net/http"

	"pulseguard/collector/internal/store"
)

type thresholdValue struct {
	Warning float64 `json:"warning"`
	Error   float64 `json:"error"`
}

type thresholdsResponse struct {
	HostID     string                    `json:"host_id"`
	Thresholds map[string]thresholdValue `json:"thresholds"`
}

func (s *Server) handleGetThresholds(w http.ResponseWriter, r *http.Request) {
	hostID := r.PathValue("host_id")

	data, err := s.store.GetThresholds(hostID)
	if err != nil {
		http.Error(w, "esikler alinamadi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := thresholdsResponse{HostID: hostID, Thresholds: make(map[string]thresholdValue)}
	for checkType, t := range data {
		resp.Thresholds[checkType] = thresholdValue{Warning: t.Warning, Error: t.Error}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleSetThresholds(w http.ResponseWriter, r *http.Request) {
	hostID := r.PathValue("host_id")

	var body struct {
		Thresholds map[string]thresholdValue `json:"thresholds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "gecersiz istek govdesi", http.StatusBadRequest)
		return
	}

	for checkType, t := range body.Thresholds {
		if t.Error <= t.Warning {
			http.Error(w, "error esigi warning esiginden buyuk olmali: "+checkType, http.StatusBadRequest)
			return
		}
		err := s.store.SetThreshold(hostID, checkType, store.ThresholdPair{Warning: t.Warning, Error: t.Error})
		if err != nil {
			http.Error(w, "esik kaydedilemedi: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
