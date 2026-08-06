# PulseGuard

Hafif filo izleme & olay toplama sistemi — config-driven bir ajan, tek-yönlü & imzalı veri akışı ve web dashboard ile uçtan uca bir sistem.

## Bileşenler

- **agent** — Go. CPU/RAM/Disk izler, olay üretir, imzalı batch olarak Collector'a gönderir.
- **collector** — Go. Batch'leri alır, imza/bütünlük doğrular, SQLite'a kaydeder, REST API sunar.
- **dashboard** — React + TypeScript. Filo görünümü, host detayı, CSV rapor indirme.
- **pulsectl** — Go CLI. `status` ve `config validate` komutları.
- **shared** — Go. Event/Batch şeması ve HMAC imzalama, tüm bileşenler arasında paylaşılan sözleşme.

## Kurulum

Gereksinimler: Go 1.22+, Node.js 20+.

### 1. Collector'ı başlat
cd collector
go run ./cmd/collector


`:8080` portunda dinlemeye başlar, `collector.db` adında bir SQLite dosyası oluşturur.

### 2. Agent'ı çalıştır

cd agent
go run ./cmd/agent


`config.example.yaml`'ı okur, CPU/RAM/Disk'i ölçer, sonuçları imzalı bir batch olarak Collector'a gönderir.

### 3. Dashboard'ı başlat

cd dashboard
npm install
npm run dev


`http://localhost:5173` adresinde açılır.

### 4. pulsectl kullanımı

cd pulsectl
go run ./cmd/pulsectl status
go run ./cmd/pulsectl config validate


## Mimari

Agent, dışarıdan hiçbir bağlantı kabul etmez — sadece Collector'a doğru bağlantı açar (tek yönlü, "diode" ilhamı). Her batch HMAC-SHA256 ile imzalanır, Collector bütünlüğü doğrulamadan hiçbir veriyi kabul etmez.

Detaylı mimari kararlar ve gerekçeleri için `docs/` klasörüne bakın.

## Ortam değişkenleri

- `PULSEGUARD_SHARED_SECRET` — Collector'ın imza doğrulamak için kullandığı gizli anahtar.
- `PULSEGUARD_AGENT_TOKEN` — Agent'ın config'inde `${PULSEGUARD_AGENT_TOKEN}` olarak referans verilen, imzalama için kullanılan anahtar (Collector'daki ile aynı olmalı).