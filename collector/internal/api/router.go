package api

import (
	"net/http"

	"pulseguard/collector/internal/store"
)

type Server struct { // api sunucusunun ihtiyaç duyduğu paylaşılan kaynakları taşıyan struct
	store  *store.Store
	secret []byte
	//HMAC doğrulama için gereken gizli anahtar , main.go'da env'den okunup buraya verilecek
}

func NewServer(s *store.Store, secret []byte) *Server { //server struct'ını oluşturan constructor fonksiyon
	return &Server{store: s, secret: secret}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux() //GO'nun standart router'ı
	mux.HandleFunc("POST /api/v1/batches", s.handleCreateBatch)
	mux.HandleFunc("GET /api/v1/hosts", s.handleListHosts)
	mux.HandleFunc("GET /api/v1/events", s.handleListEvents)
	mux.HandleFunc("GET /api/v1/thresholds/{host_id}", s.handleGetThresholds)
	mux.HandleFunc("PUT /api/v1/thresholds/{host_id}", s.handleSetThresholds)
	return withCORS(mux)
}
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST,PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
