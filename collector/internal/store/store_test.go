package store

import (
	"path/filepath"
	"testing"

	"pulseguard/shared/schema"
)

func TestSaveBatch_AndBatchExists(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open hata verdi: %v", err)
	}
	defer s.Close()
	b := schema.Batch{
		BatchID:      "batch_1",
		HostID:       "web-03",
		AgentVersion: "0.1.0",
		CreatedAt:    "2026-07-16T09:14:05Z",
		Events: []schema.Event{
			{EventID: "evt_1", HostID: "web-03", CheckType: "disk", Level: schema.LevelError, Message: "disk full", Timestamp: "2026-07-16T09:14:03Z", SchemaVersion: 1},
		},
	}

	exists, err := s.BatchExists("batch_1")
	if err != nil {
		t.Fatalf("BatchExists hata verdi: %v", err)
	}
	if exists {
		t.Error("henuz kaydedilmemis batch icin exists=true donmemeliydi")
	} //henüz hiçbir şey kaydedilmediği için false dönmeli , true dönerse exists sorgusunun mantığında sorun vardır

	if err := s.SaveBatch(b); err != nil {
		t.Fatalf("SaveBatch hata verdi: %v", err)
	}

	exists, err = s.BatchExists("batch_1")
	if err != nil {
		t.Fatalf("BatchExists hata verdi: %v", err)
	}
	if !exists {
		t.Error("kaydedilmis batch icin exists=true donmeliydi")
	}
}

func TestSaveBatch_DuplicateBatchIDFails(t *testing.T) { //idempotency kontrolünün veritabanı seviyesinde gerçekten çalıştığını kontrol ediyoruz
	dbPath := filepath.Join(t.TempDir(), "test.db")
	s, _ := Open(dbPath)
	defer s.Close()

	b := schema.Batch{
		BatchID:   "batch_1",
		HostID:    "web-03",
		CreatedAt: "2026-07-16T09:14:05Z",
	}

	if err := s.SaveBatch(b); err != nil {
		t.Fatalf("ilk SaveBatch hata verdi: %v", err)
	} //ilk kayıt başarıyla geçmeli

	err := s.SaveBatch(b)
	if err == nil {
		t.Error("ayni BatchID ikinci kez kaydedilmeye calisilinca hata donmeliydi, ama donmedi")
	} //eğer err nil ise bu primary key kısıtlamasının çalışmadığı ya da savebatch'in bu hatayı görmezden geldiği anlamına gelir. Testin geçmesi primary key constraintinin gerçek bir güvenlik katmanı olarak işlev gördüğünü kanıtlar
}
