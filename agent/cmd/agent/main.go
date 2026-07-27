package main // bu bir kütüphane paketi değil, çalıştırılabilir ve fonksiyon bulundurması gereken bir program. mutlaka bir func main() olması gerekir, bu da programın başlangıç noktası olur

import (
	"fmt"
	"os"

	"pulseguard/agent/internal/checks"
	"pulseguard/agent/internal/config"
) //check ve config paketleri import edildi

func main() { //programın çalışmaya başladığı ilk fonksiyon, program çalıştırıldığında otomatik olarak bu fonksiyon çağırılır
	cfg, err := config.Load("config.example.yaml") //daha önce yazılan load fonksiyonunu çağırır verilen dosyayı okuyup struct'a çevirecek
	if err != nil {
		fmt.Fprintf(os.Stderr, "config yuklenemedi: %v\n", err)
		os.Exit(1)
	} //hata varsa standart hata akışına yazılıyor ve program 1(başarısız) koduyla sonlanıyor

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
	}
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
