package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"pulseguard/shared/schema"
)

type Sender struct {
	collectorURL string
	httpClient   *http.Client
	backoff      []time.Duration
}

func New(collectorURL string, backoffSeconds []int) *Sender {
	backoff := make([]time.Duration, len(backoffSeconds))
	for i, sec := range backoffSeconds {
		backoff[i] = time.Duration(sec) * time.Second
	}
	return &Sender{
		collectorURL: collectorURL,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		backoff:      backoff,
	} //agent sonsuza kadar beklemesin diye 10 saniyelik bir bekleme süresi tanımlıyoruz, http isteği 10 saniyeden uzun sürerse otomatik iptal edilecek
}

func (s *Sender) Send(b schema.Batch) error { //batch'i retry mantığıyla collector'a gönderir
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("batch serialize edilemedi: %w", err)
	}

	waits := append([]time.Duration{0}, s.backoff...) //bekleme sürelerinin bşına 0 eklenir ki ilk deneme hiç beklemeden yapılsın, sadece başarısız olursa sonraki denemeler backoff sürelerine göre bekler

	var lastErr error
	for attempt, wait := range waits {
		if wait > 0 {
			time.Sleep(wait)
		}

		resp, err := s.httpClient.Post(s.collectorURL+"/api/v1/batches", "application/json", bytes.NewReader(data)) //gerçek bir http post isteği gönderiyoruz
		if err != nil {
			lastErr = fmt.Errorf("istek gonderilemedi (deneme %d): %w", attempt+1, err)
			continue
			//ağ hatası olursa hatayı kaydedip bir sonraki denemeye geçiyoruz
		}
		resp.Body.Close() //bağlantı kaynaklarının sızmaması için http cevabının gövdesini kapatıyoruz

		if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusConflict {
			return nil
			//201 veya 409 dönerse başarılı sayılır collector veriyi biliyor anlamına gelir
		}
		lastErr = fmt.Errorf("collector hata dondu (deneme %d): HTTP %d", attempt+1, resp.StatusCode)
		//400,422,500 gibi farklı bir durum kodu dönerse bu hata olarak kaydedilir ve bir sonraki denemeye geçilir
	}

	return lastErr //tüm denemeler tükendiğinde en son hatayı döndürürüz bu da agent ana döngüsüne bu batch gönderilemedi kuyrukta beklesin sinyali gönderir
}
