# PulseGuard — Mimari Tasarım Notları

## 1. Genel Bakış

PulseGuard, sunucularda çalışan bir izleme ajanı (Agent), bu verileri toplayan merkezi bir servis (Collector) ve verileri gösteren bir web panelinden (Dashboard) oluşan bir filo izleme sistemidir. Sistem dört bağımsız parçaya bölünmüştür: Agent, Collector, Dashboard, pulsectl (CLI). Her parçanın tek bir sorumluluğu vardır ve parçalar birbirleriyle yalnızca tanımlı arayüzler (JSON/REST) üzerinden konuşur.

## 2. Neden Ayrı Servisler (Monolitik Değil)?

Sistemi tek bir parça yerine ayrı servislere böldük çünkü:
- Her parça bağımsız geliştirilip test edilebiliyor (örneğin `checks` paketi, hiçbir ağ bağlantısı olmadan test edilebildi).
- Güvenlik sınırları netleşiyor: Agent'ın çöküşü Collector'ı etkilemez, Collector'ın çöküşü Agent'ın veri kaybetmesine yol açmaz (yerel kuyruk sayesinde).
- Farklı teknolojiler (Go / React+TS) doğal olarak ayrı süreçler gerektiriyor.

Dezavantajı: dağıtık bir sistemde ağ hataları, sıralama sorunları gibi ekstra karmaşıklıklar yönetilmesi gerekiyor — bunu retry/backoff ve idempotency ile çözdük.

## 3. Tek Yönlülük İlkesi (Diode İlhamı)

Agent, dışarıdan hiçbir bağlantı kabul etmez — dinleyen bir portu yoktur, sadece kendisi Collector'a doğru bağlantı açar. Bu tasarım, DataFlowX'in `diode` ürününden ilham alıyor: iki izole ağ arasında veriyi yalnızca tek yöne akıtan bir güvenlik diyotu. Bizim projemizde bunun karşılığı, Collector ele geçirilse bile saldırganın Agent'ın çalıştığı makineye geri sızamamasıdır.

## 4. Veri Bütünlüğü — İmza

Her batch, Agent tarafından HMAC-SHA256 ile imzalanır (`shared/schema/signature.go`). Agent ve Collector, aynı `canonicalPayload` fonksiyonunu kullanarak imzayı hesaplar — bu, "Agent nasıl imzalıyorsa Collector da tam öyle doğruluyor" garantisini sağlıyor. Doğrulamada `hmac.Equal` kullanılıyor (düz `==` değil) çünkü bu, zamanlama saldırılarına (timing attack) karşı sabit-zamanlı bir karşılaştırma sağlıyor.

**Neden HMAC (asimetrik imza değil)?** HMAC, Go'nun standart kütüphanesinde hazır, kurulumu basit ve "veri bütünlüğü" dersini tam veriyor. Asimetrik imza (ed25519) daha güçlü bir güvenlik modeli sunar (private key hiç ağda dolaşmaz) ama key management karmaşıklığı ekliyor — bu yüzden brief'in "ileri hedef" (F13) olarak bıraktığı bu özelliği projenin çekirdeğine dahil etmedik.

## 5. Dayanıklılık — Yerel Kuyruk ve Retry

Agent, ürettiği her event'i önce yerel bir kuyruğa (`agent/internal/queue`, JSONL formatında, diske yazılan) ekliyor. Collector'a gönderim başarısız olursa, `retry_backoff_seconds` listesine göre kademeli artan beklemelerle tekrar deniyor. Gönderim başarılı olursa (ya da Collector "zaten işlendi" derse — 409), event'ler kuyruktan siliniyor; başarısız olursa kuyrukta kalıp bir sonraki çalıştırmada tekrar denenmeye devam ediyor.

**Neden diske yazan bir kuyruk (bellekte değil)?** Agent süreci yeniden başlarsa (çökme, yeniden başlatma), bellekteki veri kaybolurdu. Diske yazmak, bu senaryoda bile event'lerin korunmasını sağlıyor.

## 6. Idempotency — Aynı Batch İki Kez Gelirse

Collector, her batch'i `batch_id` ile SQLite'ta `PRIMARY KEY` kısıtlamalı bir tabloya kaydediyor. Aynı `batch_id` ikinci kez gelirse, veritabanı bunu otomatik reddediyor, Collector bu durumu `409 Conflict` olarak Agent'a bildiriyor. Agent, hem `201` (yeni kabul) hem `409` (zaten var) durumlarını "başarılı" sayıp kuyruğu temizliyor — çünkü her iki durumda da Collector artık bu veriyi biliyor.

## 7. Neden SQLite (PostgreSQL Değil)?

SQLite, tek bir dosya olarak çalışıyor, ayrı bir sunucu/kurulum gerektirmiyor — staj ölçeğinde (birkaç host, orta hacimli event) yeterli. `modernc.org/sqlite` kütüphanesini seçtik çünkü CGO gerektirmiyor (saf Go), bu da cross-compile'ı (özellikle Docker'a geçerken) kolaylaştırıyor, hiçbir C derleyicisine ihtiyaç duymuyor.

PostgreSQL'e geçiş ihtiyacı; çoklu-yazıcı eşzamanlılığı, çok büyük veri hacmi ya da yüksek eşzamanlı sorgu gerektiğinde doğar — bu proje ölçeğinde bu ihtiyaç yok.

## 8. Push vs Pull İzleme Modeli

Agent, verisini kendisi Collector'a **push** ediyor (Collector'ın Agent'a bağlanıp veri "çekmesi" değil). Bunun sebebi doğrudan tek-yönlülük ilkesiyle bağlantılı: pull modelinde Collector'ın her Agent'a bağlanabilmesi gerekir, bu da her Agent'ın bir port açık tutmasını gerektirir — tam tersi bir güvenlik profili. Push modeli, Agent'ın hiçbir zaman dışarıdan erişilebilir olmamasını sağlıyor.

## 9. Event vs Metrik

Sistemde hem "olay" (event — bir eşik aşıldığında üretilen, seviyeli bir kayıt) hem de "metrik" (event içindeki `Metric` alanı — ham sayısal ölçüm) birlikte tutuluyor. Sadece metrik toplasaydık, "bu değer normal mi anormal mi" yorumunu her sorguda yeniden yapmamız gerekirdi; sadece event toplasaydık, ham sayısal veriyi (örn. grafik çizmek için) kaybederdik. İkisini birlikte tutmak, hem alarm mantığını hem de görselleştirmeyi tek bir veri modelinden besleyebilmemizi sağlıyor.

## 10. Neden `shared` Ayrı Bir Modül?

Agent, Collector ve pulsectl arasında paylaşılan tek şey Event/Batch şeması ve imza mantığıdır. Bunu ayrı bir `shared` modülünde tutmak, bu "sözleşmenin" birden fazla yerde kopyalanıp zamanla birbirinden sapmasını (örneğin bir tarafta yeni alan eklenip diğerinde unutulmasını) derleyici seviyesinde engelliyor. `go.work`, bu ayrı modüllerin (`agent`, `collector`, `shared`, `pulsectl`) yerel geliştirmede birbirine bağlanmasını sağlıyor; her biri yine de kendi `go.mod`'una sahip, bağımsız derlenebilir kalıyor.

## 11. Neden Polling (WebSocket Değil)?

Dashboard, Collector'dan verileri 5 saniyede bir `fetch` ile çekiyor (polling), WebSocket kullanmıyoruz. Bu projenin veri hacmi ve gerçek-zamanlılık ihtiyacı, WebSocket'in getirdiği karmaşıklığı (bağlantı yönetimi, yeniden bağlanma mantığı) haklı çıkarmıyor — birkaç saniyelik gecikme, bu ölçekte kabul edilebilir.

## 12. Sırlar ve Yapılandırma

`PULSEGUARD_SHARED_SECRET` (Collector) ve `PULSEGUARD_AGENT_TOKEN` (Agent) hiçbir zaman koda ya da config dosyasına açık yazılmıyor — ortam değişkenlerinden okunuyor. Agent'ın config dosyasında `${PULSEGUARD_AGENT_TOKEN}` gibi bir şablon görülürse, bu çalışma anında gerçek env değerine çözülüyor (`agent/internal/config/resolveEnv`).

## 13. Kimlik Üretimi

Event ve Batch ID'leri, `crypto/rand` ile üretilen 16 baytlık rastgele değerlerin hex string'e çevrilmesiyle oluşturuluyor. ULID gibi zaman-sıralı bir kimlik yerine bu basit yöntemi seçtik çünkü zaman bilgisi zaten her Event'in `Timestamp` alanında ayrıca tutuluyor — ID'nin de zaman taşıması gereksiz bir tekrar olurdu. Bu tercih, ekstra bir bağımlılık eklemekten kaçınma felsefesiyle de örtüşüyor.

