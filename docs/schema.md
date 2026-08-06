# PulseGuard — Veri Şeması

Bu doküman, Agent ile Collector arasındaki veri sözleşmesini ve Collector'ın Dashboard'a sunduğu REST API'nin şeklini açıklar. Go karşılıkları `shared/schema/` altında tanımlıdır.

## 1. Event (Olay)

Bir check'in ürettiği tek bir kayıt.

```json
{
  "event_id": "063e32f399b2daa4305ef7a958e3b8e0",
  "host_id": "local-test-host",
  "check_type": "disk",
  "level": "WARNING",
  "message": "Disk kullanimi 88.0%, uyari esigi 85.0% asildi (C:)",
  "metric": {
    "name": "disk_used_percent",
    "value": 88.0,
    "unit": "percent"
  },
  "timestamp": "2026-08-04T09:01:59Z",
  "schema_version": 1
}
```

| Alan | Tip | Açıklama |
|---|---|---|
| `event_id` | string | Benzersiz kimlik (16 byte rastgele, hex) |
| `host_id` | string | Olayı üreten sunucunun kimliği |
| `check_type` | string | `"cpu"`, `"ram"`, `"disk"` |
| `level` | string | `"INFO"` \| `"WARNING"` \| `"ERROR"` |
| `message` | string | İnsan-okur açıklama |
| `metric` | object | Ham ölçüm (opsiyonel — `omitempty`) |
| `timestamp` | string | RFC3339 formatında, UTC |
| `schema_version` | int | Şema sürüm numarası (şu an sabit `1`) |

**Neden `metric` opsiyonel?** İleride eklenebilecek bazı check tiplerinin (örn. bir HTTP endpoint'inin sadece "ayakta/kapalı" bilgisini veren bir check) sayısal bir ölçümü olmayabilir. `omitempty` ile bu durumda alan JSON çıktısında hiç görünmez, boş/sıfır değerli bir obje yerine.

## 2. Batch (Toplu Gönderim Zarfı)

Agent, event'leri tek tek değil, toplu (batch) halinde gönderir.

```json
{
  "batch_id": "a1b2c3d4e5f6...",
  "host_id": "local-test-host",
  "agent_version": "0.1.0",
  "created_at": "2026-08-04T09:01:59Z",
  "events": [ "... Event objeleri ..." ],
  "signature": {
    "algorithm": "hmac-sha256",
    "value": "b17a3c...9f2e"
  }
}
```

| Alan | Tip | Açıklama |
|---|---|---|
| `batch_id` | string | Benzersiz kimlik, idempotency kontrolünde kullanılır |
| `host_id` | string | Batch'i gönderen sunucu |
| `agent_version` | string | Gönderen Agent'ın sürümü |
| `created_at` | string | RFC3339, batch'in oluşturulduğu an |
| `events` | array | Event listesi |
| `signature` | object | HMAC-SHA256 imzası |

**İmza nasıl hesaplanır:** `batch_id`, `host_id`, `agent_version`, `created_at` ve `events` alanları, sabit sırayla JSON'a çevrilip (`canonicalPayload`), bu byte dizisi HMAC-SHA256 ile imzalanır. `signature` alanının kendisi imza hesaplamasına dahil edilmez.

## 3. REST API

### `POST /api/v1/batches`

Agent'tan gelen batch'i kabul eder.

| Durum Kodu | Anlamı |
|---|---|
| `201 Created` | Batch başarıyla kaydedildi |
| `400 Bad Request` | JSON gövdesi bozuk/parse edilemedi |
| `409 Conflict` | Bu `batch_id` daha önce işlenmiş (idempotency) |
| `422 Unprocessable Entity` | İmza doğrulaması başarısız |
| `500 Internal Server Error` | Beklenmeyen veritabanı hatası |

### `GET /api/v1/hosts`

Filo özetini döner.

```json
{
  "hosts": [
    { "host_id": "web-03", "status": "ok", "last_seen": "2026-08-04T09:01:59Z" }
  ]
}
```

`status`, o host'un en son event'inin seviyesine göre hesaplanır: `ERROR` → `"error"`, `WARNING` → `"warning"`, aksi halde `"ok"`; hiç event yoksa `"unknown"`.

### `GET /api/v1/events?host_id=<id>&limit=<n>`

Event geçmişini döner (en yeniden eskiye sıralı). `host_id` verilmezse tüm host'ların event'leri gelir. `limit` verilmezse varsayılan `50`.

```json
{
  "events": [ "... Event objeleri ..." ]
}
```

## 4. Agent Config (YAML)

```yaml
collector:
  url: "http://localhost:8080"
  auth_token: "${PULSEGUARD_AGENT_TOKEN}"
  batch_interval_seconds: 30
  retry_backoff_seconds: [5, 15, 60, 300]

host:
  id: "local-test-host"
  tags: ["dev", "test"]

checks:
  - type: cpu
    interval_seconds: 15
    threshold_warning: 75
    threshold_error: 90
  - type: disk
    interval_seconds: 60
    mount: "C:"
    threshold_warning: 85
    threshold_error: 92

queue:
  path: "./queue.db"
  max_size_mb: 200
```

`auth_token` alanındaki `${...}` kalıbı, Agent çalışma anında gerçek ortam değişkeni değeriyle değiştirilir — sır asla dosyada açık yazılmaz.