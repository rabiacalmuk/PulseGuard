package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"pulseguard/collector/internal/verify"
	"pulseguard/shared/schema"
)

func (s *Server) handleCreateBatch(w http.ResponseWriter, r *http.Request) {
	var b schema.Batch
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "gecersiz istek govdesi", http.StatusBadRequest)
		return
	}

	err := verify.Batch(s.secret, s.store, b)
	if errors.Is(err, verify.ErrDuplicateBatch) {
		http.Error(w, "batch daha once islenmis", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "dogrulama basarisiz: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := s.store.SaveBatch(b); err != nil {
		http.Error(w, "kaydetme basarisiz: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
