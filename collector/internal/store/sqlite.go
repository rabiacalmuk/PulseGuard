package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"pulseguard/shared/schema"
)

type Store struct { //Collector'ın  veritabanı bağlantısını taşıyan struct
	db *sql.DB
}

func Open(path string) (*Store, error) { //veritabanını açan ve gerekli tabloları oluşturan fonksiyon
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("veritabani acilamadi: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migration basarisiz: %w", err)
	} //aşağıdaki migrate fonksiyonunu açıp gerekli tabloların var olduğundan emin olunur

	return &Store{db: db}, nil
}
func (s *Store) Close() error {
	return s.db.Close()
}

func migrate(db *sql.DB) error { //veritabanının yapısını oluşturan, güncelleyen işlem
	schema := `
	CREATE TABLE IF NOT EXISTS batches (
		batch_id TEXT PRIMARY KEY,
		host_id TEXT NOT NULL,
		created_at TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS events (
		event_id TEXT PRIMARY KEY,
		batch_id TEXT NOT NULL,
		host_id TEXT NOT NULL,
		check_type TEXT NOT NULL,
		level TEXT NOT NULL,
		message TEXT NOT NULL,
		metric_name TEXT,
		metric_value REAL,
		metric_unit TEXT,
		timestamp TEXT NOT NULL,
		schema_version INTEGER NOT NULL
	);
	`
	_, err := db.Exec(schema)
	return err
}

func (s *Store) BatchExists(batchID string) (bool, error) { //bu batchID daha önce işlendi mi kontrolü yapar
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM batches WHERE batch_id = ?)", batchID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("batch kontrolu basarisiz: %w", err)
	}
	return exists, nil
}

func (s *Store) SaveBatch(b schema.Batch) error { //batch'i ve içindeki tüm eventleri kalıcı olarak kaydeden fonksiyon
	tx, err := s.db.Begin() //transaction başlatılır
	if err != nil {
		return fmt.Errorf("transaction baslatilamadi: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO batches (batch_id, host_id, created_at) VALUES (?, ?, ?)",
		b.BatchID, b.HostID, b.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("batch kaydedilemedi: %w", err)
	} //önce batch satırını ekledik ki bu batchID zaten varsa baştan hata yakalayıp döndürelim

	for _, e := range b.Events {
		_, err = tx.Exec(
			`INSERT INTO events
				(event_id, batch_id, host_id, check_type, level, message, metric_name, metric_value, metric_unit, timestamp, schema_version)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			e.EventID, b.BatchID, e.HostID, e.CheckType, e.Level, e.Message,
			e.Metric.Name, e.Metric.Value, e.Metric.Unit, e.Timestamp, e.SchemaVersion,
		)
		if err != nil {
			return fmt.Errorf("event kaydedilemedi (%s): %w", e.EventID, err)
		} //batch'in tüm eventleri için döngüye girilir ve eventlerin tüm alanları ve batchID (eventin hangi batch'e ait olduğunu kaydederiz) eklenir
	}

	return tx.Commit()
	//tüm eklemeler başarılıysa tansaction onaylanır ve değişiklikler kalıcı hale gelir. Buraya kadar bir hata olsaydı deferdaki rollback tüm değişiklikleri geri alırdı

}
