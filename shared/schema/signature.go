package schema

import (
	"crypto/hmac"   //HMAC (Hash based message authentication code) hesaplamak için kullanılan paket
	"crypto/sha256" //HMAC'in hangi hash fonksiyonuyla çalışacağını belirtir
	"encoding/hex"  //hesaplanan imza değerini hex formatına çevirmeye kullanılır
	"encoding/json" //batch'i imza hesaplamadan önce json byte dizisine çevirmek için kullanılır
	"fmt"           //hata mesajlarını biçimlendirmek için kullanılır
)

func canonicalPayload(b Batch) ([]byte, error) { //Bu fonksiyon ile imzalanacak veri tutarlı bir byte dizisine çevrilir
	payload := struct {
		BatchID      string  `json:"batch_id"`
		HostID       string  `json:"host_id"`
		AgentVersion string  `json:"agent_version"`
		CreatedAt    string  `json:"created_at"`
		Events       []Event `json:"events"`
	}{ // Batch'in Signature alanını imzalama hesaplamasına dahil etmemek amacıyla yalnızca imzalanacak alanlardan oluşan anonim bir struct oluşturulur.
		BatchID:      b.BatchID,
		HostID:       b.HostID,
		AgentVersion: b.AgentVersion,
		CreatedAt:    b.CreatedAt,
		Events:       b.Events,
	}
	return json.Marshal(payload) //struct'ı json byte dizisine çevirir ve döndürür
}

func Sign(secret []byte, b Batch) (Signature, error) { //secret:agent ve collector'ın paylaştığı gizli anahtar, b:imzalanacak batch
	payload, err := canonicalPayload(b)
	if err != nil { // err nil değilse hata var demektir. %w, orijinal hatayı "sarmalayarak" (wrap) saklar ileride bu hatanın kök nedenini errors.Is/As ile takip edebiliriz.
		return Signature{}, fmt.Errorf("payload oluşturulamadı: %w", err)
	}

	mac := hmac.New(sha256.New, secret) //SHA-256 tabanlı bir HMAC hesaplayıcı nesnesi oluşturur, secret ile başlatılır
	mac.Write(payload)                  //cononical Json byte dizisini HMAC hesaplayıcıya yazar, imzza bu veri üzerindeen hesaplanacak
	sum := mac.Sum(nil)                 // hesaplamayı tamamlayıp imza değerini byte dizisi olarak döndürür

	return Signature{
		Algorithm: "hmac-sha256",           //kullandığımız algoritmayı belirtiyoruz, collector bu bilgiye bakıp doğru doğrulama yolunu seçer
		Value:     hex.EncodeToString(sum), //binary imza byte'larını okunabilir bir hex stringe çevirir
	}, nil
}

func Verify(secret []byte, b Batch) error { //collector tarafında çağırılacak fonksiyon, gelen batch'ın imzasının geçerli olup olmadığını kontrol eder
	expected, err := Sign(secret, b) //doğrulama işlemini imzayı çözerek değil aynı hesaplamayı tekrarlayarak yapar
	if err != nil {
		return fmt.Errorf("doğrulama için imza hesaplanamadı: %w", err)
	}

	if b.Signature.Algorithm != expected.Algorithm {
		return fmt.Errorf("desteklenmeyen imza algoritması: %s", b.Signature.Algorithm)
	}

	if !hmac.Equal([]byte(b.Signature.Value), []byte(expected.Value)) {
		return fmt.Errorf("imza doğrulaması başarısız") //hmac.Equal, zamanlama saldırılarına karşı güvenli bir şekilde iki byte dizisini karşılaştırır
	}

	return nil
}
