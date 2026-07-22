package checks

import (
	"fmt"
	"time" //zaman ölçümü yapmak için bu paketi kullanıyoruz

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUCheck struct { //check'in ihtiyaç duyduğu ayarları taşıyan struct
	ThresholdWarn  float64
	ThresholdError float64
}

func (c CPUCheck) Name() string {
	return "cpu"
}

func (c CPUCheck) Run() (CheckResult, error) {
	percentages, err := cpu.Percent(time.Second, false) //gopsutil'in cpu.Percent fonksiyonu çağrılır, bu OS'den gerçek CPU kullanım bilgisini okur, time.Second ile 1 saniye boyunca ölçüm yapılması sağlanır, false ile tüm çekirdeklerin ortalaması alınır
	if err != nil {
		return CheckResult{}, fmt.Errorf("cpu kullanimi okunamadi: %w", err)
	}
	if len(percentages) == 0 {
		return CheckResult{}, fmt.Errorf("cpu kullanimi okunamadi: bos sonuc")
	} //gopsutil'in cpu.Percent fonksiyonu boş bir slice döndürürse hata mesajı döndürülür

	percent := percentages[0] //tüm çekirdeklerin ortalaması alındığı için slice'ın ilk elemanı kullanılır

	level := DetermineLevel(percent, c.ThresholdWarn, c.ThresholdError)

	var message string
	switch level {
	case LevelError:
		message = fmt.Sprintf("CPU kullanimi %.1f%%, hata esigi %.1f%% asildi ", percent, c.ThresholdError)
	case LevelWarning:
		message = fmt.Sprintf("CPU kullanimi %.1f%%, uyari esigi %.1f%% asildi ", percent, c.ThresholdWarn)
	default:
		message = fmt.Sprintf("CPU kullanimi %.1f%% ", percent)
	}

	return CheckResult{
		CheckType:   "cpu",
		Level:       level,
		Message:     message,
		MetricName:  "cpu_used_percent",
		MetricValue: percent,
		MetricUnit:  "percent",
	}, nil
}
