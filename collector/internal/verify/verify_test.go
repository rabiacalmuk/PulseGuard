package verify

import (
	"errors" //errors.IS ile belirli bir sentinel error olup olmadığını kontrol etmek için import edildi
	"path/filepath"
	"testing"

	"pulseguard/collector/internal/store"
	"pulseguard/shared/schema"
)

func newTestStore(t *testing.T) *store.Store {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	s, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("store acilamadi: %v", err)
	}
	t.Cleanup(func() { s.Close() }) //s.Close()'u her seferinde elle yazmak yerine bu yardımcı fonksiyon içinde bir kez tanımlanır ve her çağıran test otomatik faydalanır bundan
	return s
}

func signedBatch(t *testing.T, secret []byte, batchID string) schema.Batch { //geçerli şekilde imzalanmış bir batch oluşturma işlemini tek yerde topladık
	b := schema.Batch{
		BatchID:   batchID,
		HostID:    "web-03",
		CreatedAt: "2026-07-16T09:14:05Z",
	}
	sig, err := schema.Sign(secret, b)
	if err != nil {
		t.Fatalf("Sign hata verdi: %v", err)
	}
	b.Signature = sig
	return b
}

func TestBatch_ValidAndNew_Succeeds(t *testing.T) {
	secret := []byte("test-secret")
	s := newTestStore(t)

	b := signedBatch(t, secret, "batch_1")

	if err := Batch(secret, s, b); err != nil {
		t.Errorf("gecerli, yeni bir batch icin hata donmemeliydi: %v", err)
	}
}

func TestBatch_WrongSignature_Fails(t *testing.T) { //yanlış anahtarla doğrulama yapmayı dener başarısız olmalı
	s := newTestStore(t)

	b := signedBatch(t, []byte("dogru-anahtar"), "batch_1") //doğru anahtarla imzalanıyor

	err := Batch([]byte("yanlis-anahtar"), s, b) //farklı bir anahtarla doğrulamayı deniyoruz
	if err == nil {
		t.Error("yanlis anahtarla dogrulama basarili olmamaliydi")
	}
	if errors.Is(err, ErrDuplicateBatch) {
		t.Error("bu hata ErrDuplicateBatch olmamaliydi, imza hatasi olmaliydi")
	} //ErrDuplicatedBatch dönseydi bu Batch fonksiyonunun sıralamasının bozulduğunu gösterirdi
}

func TestBatch_Duplicate_ReturnsErrDuplicateBatch(t *testing.T) { //aynı batch iki kez batch fonksiyonuna verilirse ikinci seferde ErrDuplicatedBatch dönmeli
	secret := []byte("test-secret")
	s := newTestStore(t)

	b := signedBatch(t, secret, "batch_1")

	if err := Batch(secret, s, b); err != nil {
		t.Fatalf("ilk dogrulama basarisiz olmamaliydi: %v", err)
	} //henüz veritabanında olmadığı için başarılı olmalı
	if err := s.SaveBatch(b); err != nil {
		t.Fatalf("SaveBatch basarisiz oldu: %v", err)
	} //veritabanına kaydediliyor artık

	err := Batch(secret, s, b) //aynı batch ile ikinci kez doğrulama denenir
	if !errors.Is(err, ErrDuplicateBatch) {
		t.Errorf("ErrDuplicateBatch bekleniyordu, gelen: %v", err)
	} //errors.Is ile ErrDuplicatedBatch olup olmadığını ya da onu sarmalayıp sarmalamadığını kontrol eder
}
