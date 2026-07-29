package verify

import (
	"fmt"

	"pulseguard/collector/internal/store" //BatchExists'i çağırabilmek için store paketini import ettik
	"pulseguard/shared/schema"
)

var ErrDuplicateBatch = fmt.Errorf("batch daha once islenmis")

func Batch(secret []byte, s *store.Store, b schema.Batch) error {
	if err := schema.Verify(secret, b); err != nil {
		return fmt.Errorf("imza dogrulama basarisiz: %w", err)
	} //imza doğrulaması yapılıyor eğer başarısızsa idempotency kontrolüne gidilmiyor

	exists, err := s.BatchExists(b.BatchID)
	if err != nil {
		return fmt.Errorf("batch kontrolu basarisiz: %w", err)
	}
	if exists {
		return ErrDuplicateBatch
	} //imza geçerliyse duplication kontrolü yapılır

	return nil
} //iki kontrol de doğruysa batch güvenle kaydedilebilir
