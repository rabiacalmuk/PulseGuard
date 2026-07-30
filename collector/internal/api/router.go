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

func (s *Server) Router() *http.ServeMux { //URL'leri handlerlara bağlayan fonksiyon
	mux := http.NewServeMux() //GO'nun standart router'ı
	mux.HandleFunc("POST /api/v1/batches", s.handleCreateBatch)
	return mux
}
