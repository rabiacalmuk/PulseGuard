package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type eventsResponse struct {
	Events []interface{} `json:"events"`
}

func (s *Server) handleListEvents(w http.ResponseWriter, r *http.Request) {
	hostID := r.URL.Query().Get("host_id")

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	events, err := s.store.ListEvents(hostID, limit)
	if err != nil {
		http.Error(w, "event listesi alinamadi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"events": events})
}
