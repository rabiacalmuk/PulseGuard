package main // bu bir kütüphane paketi değil, çalıştırılabilir ve fonksiyon bulundurması gereken bir program. mutlaka bir func main() olması gerekir, bu da programın başlangıç noktası olur

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"pulseguard/agent/internal/checks"
	"pulseguard/agent/internal/config"
	"pulseguard/agent/internal/queue"
	"pulseguard/agent/internal/sender"
) // paketler import edildi

func main() { //programın çalışmaya başladığı ilk fonksiyon, program çalıştırıldığında otomatik olarak bu fonksiyon çağırılır
	cfg, err := config.Load("config.example.yaml") //daha önce yazılan load fonksiyonunu çağırır verilen dosyayı okuyup struct'a çevirecek
	if err != nil {
		fmt.Fprintf(os.Stderr, "config yuklenemedi: %v\n", err)
		os.Exit(1)
	} //hata varsa standart hata akışına yazılıyor ve program 1(başarısız) koduyla sonlanıyor

	q, err := queue.Open(cfg.Queue.Path) //config'deki queue.path ilekuyruğu açıyoruz
	if err != nil {
		fmt.Fprintf(os.Stderr, "kuyruk acilamadi: %v\n", err)
		os.Exit(1)
	}

	for _, checkCfg := range cfg.Checks { //cfg.Checks listesindeki her eleman için döngüye girilir
		c, err := buildCheck(checkCfg) //yaml'dan gelen ham checkConfig'i gerçek bir checks.Check nesnesine çevirir
		if err != nil {
			fmt.Fprintf(os.Stderr, "check olusturulamadi (%s): %v\n", checkCfg.Type, err)
			continue
		}

		result, err := c.Run() //gopsutil ile gerçek ölçümü yapan kısım
		if err != nil {
			fmt.Fprintf(os.Stderr, "check calistirilamadi (%s): %v\n", c.Name(), err)
			continue
		}

		fmt.Printf("[%s] %s: %s\n", result.Level, result.CheckType, result.Message)

		event := sender.ToEvent(cfg.Host.ID, result) //ham checkResult gerçek bir schemaEvent'e dönüştürüyoruz
		if err := q.Enqueue(event); err != nil {
			fmt.Fprintf(os.Stderr, "kuyruga eklenemedi: %v\n", err) //çevrilen event kuyruğa ekleniyor
		}
	}
	events, err := q.All()
	if err != nil {
		fmt.Fprintf(os.Stderr, "kuyruk okunamadi: %v\n", err)
		os.Exit(1)
		//tüm döngü bittikten sonra kuyrukta biriken tüm her şeyi okuyor(sadece bu çalıştırmadan değil önceki başarısız denemelerden kalan eventler de okunuyor)
	}
	if len(events) == 0 {
		fmt.Println("gonderilecek event yok")
		return
	} //gönderilecek event yoksa boş batch göndermemek için çıkıyoruz direkt

	secret := []byte(cfg.Collector.AuthToken) //configteki auth_token imzalama için gereken secret olarak kullanılıyor
	batch, err := sender.BuildBatch(secret, cfg.Host.ID, "0.1.0", events)
	if err != nil {
		fmt.Fprintf(os.Stderr, "batch olusturulamadi: %v\n", err)
		os.Exit(1)
	} //toplanan eventler imzalı bir batch haline getirilir, "0.1.0" şu an için sabit bir agent sürüm numarası

	s := sender.New(cfg.Collector.URL, cfg.Collector.RetryBackoffSeconds)
	if err := s.Send(batch); err != nil {
		fmt.Fprintf(os.Stderr, "gonderim basarisiz, event'ler kuyrukta kaldi: %v\n", err)
		return
	} //gerçek gönderim yapılıyor eğer başarısız olursa evenler kuyrukta kalıyor ve sadece hata mesajı basarak programdan çıkıyoruz, sonraki çalıştırmada bu eventler tekrar denenecek

	ids := make([]string, len(events))
	for i, e := range events {
		ids[i] = e.EventID
	} //gönderilen eventlerin id'leri ayrı bir listede toplanıyor, ack() fonksiyonu bunu bekliyor
	if err := q.Ack(ids); err != nil {
		fmt.Fprintf(os.Stderr, "kuyruk temizlenemedi: %v\n", err)
	} //gönderim başarılıysa bu eventler kuyruktan siliniyor çünkü zaten collector'a ulaştılar

	fmt.Println("batch basariyla gonderildi")
	writeStatus(len(events))
}

type agentStatus struct {
	LastRunAt  string `json:"last_run_at"`
	EventsSent int    `json:"events_sent"`
}

func writeStatus(eventCount int) {
	status := agentStatus{
		LastRunAt:  time.Now().UTC().Format(time.RFC3339),
		EventsSent: eventCount,
	}
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile("status.json", data, 0644)

}

func buildCheck(cfg config.CheckConfig) (checks.Check, error) { //yaml'dan gelen ham CheckConfig'i alıp doğru tipte bir checks.Check döndüren bir yardımcı fonksiyondur
	switch cfg.Type {
	case "cpu":
		return checks.CPUCheck{
			ThresholdWarn:  cfg.ThresholdWarning,
			ThresholdError: cfg.ThresholdError,
		}, nil //CheckConfigdeki eşik değerlerini gerçek cpuCheck struct'ına aktarıyor
	case "ram":
		return checks.RAMCheck{
			ThresholdWarn:  cfg.ThresholdWarning,
			ThresholdError: cfg.ThresholdError,
		}, nil
	case "disk":
		return checks.DiskCheck{
			Mount:          cfg.Mount,
			ThresholdWarn:  cfg.ThresholdWarning,
			ThresholdError: cfg.ThresholdError,
		}, nil //disk için mount alanını da aktarıyor
	default:
		return nil, fmt.Errorf("bilinmeyen check tipi: %s", cfg.Type)
	}
}
