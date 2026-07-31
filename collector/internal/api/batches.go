package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"pulseguard/collector/internal/verify"
	"pulseguard/shared/schema"
)

func (s *Server) handleCreateBatch(w http.ResponseWriter, r *http.Request) {
	var b schema.Batch
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		log.Printf("istek reddedildi: gecersiz govde: %v", err)
		http.Error(w, "gecersiz istek govdesi", http.StatusBadRequest)
		return
	}

	err := verify.Batch(s.secret, s.store, b)
	if errors.Is(err, verify.ErrDuplicateBatch) {
		log.Printf("batch reddedildi (tekrar gonderim): batch_id=%s host_id=%s", b.BatchID, b.HostID)
		http.Error(w, "batch daha once islenmis", http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("batch reddedildi (dogrulama basarisiz): batch_id=%s host_id=%s hata=%v", b.BatchID, b.HostID, err)
		http.Error(w, "dogrulama basarisiz: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := s.store.SaveBatch(b); err != nil {
		log.Printf("batch kaydedilemedi: batch_id=%s hata=%v", b.BatchID, err)
		http.Error(w, "kaydetme basarisiz: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("batch kabul edildi: batch_id=%s host_id=%s event_sayisi=%d", b.BatchID, b.HostID, len(b.Events))
	w.WriteHeader(http.StatusCreated)
}
