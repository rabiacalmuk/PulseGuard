package main

import (
	"log" //otomatik olarak her satıra zaman damgası ekliyor
	"net/http"
	"os"

	"pulseguard/collector/internal/api"
	"pulseguard/collector/internal/notify"
	"pulseguard/collector/internal/store"
)

func main() {
	dbPath := "collector.db" //collector çalıştığı klasörde collector.db adında bir sqlite dosyası oluşturacak
	s, err := store.Open(dbPath)
	if err != nil {
		log.Fatalf("store acilamadi: %v", err)
	}
	defer s.Close()

	secret := []byte(os.Getenv("PULSEGUARD_SHARED_SECRET")) //sır koda hiç yazılmıyor env değişkeninden okunuyor
	if len(secret) == 0 {
		log.Println("uyari: PULSEGUARD_SHARED_SECRET tanimli degil, bos anahtar kullaniliyor")
	}

	webhookURL := os.Getenv("PULSEGUARD_WEBHOOK_URL")
	notifier := notify.New(webhookURL)
	if notifier.Enabled() {
		log.Printf("webhook bildirimleri aktif: %s", webhookURL)
	} else {
		log.Println("webhook bildirimleri kapali (PULSEGUARD_WEBHOOK_URL tanimli degil)")
	}

	server := api.NewServer(s, secret, notifier) //server'ı store ve secret ile oluşturuyoruz

	addr := ":8080"
	log.Printf("collector %s adresinde dinliyor", addr)
	if err := http.ListenAndServe(addr, server.Router()); err != nil {
		log.Fatalf("sunucu baslatilamadi: %v", err)
	} //http.ListenAndServego'nun standart http sunucusunu başlatan fonksiyon, bu satır program kapanana kadar devam eder ve glen her isteği router'a yönlendirir
}
