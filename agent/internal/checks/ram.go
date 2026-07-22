package checks

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem" //gopsutil memory alt paketi
)

type RAMCheck struct {
	ThresholdWarn  float64
	ThresholdError float64
}

func (r RAMCheck) Name() string {
	return "ram"
}

func (r RAMCheck) Run() (CheckResult, error) {
	vmem, err := mem.VirtualMemory() //gopsutil'in mem.VirtualMemory fonksiyonu çağrılır, bu OS'den gerçek RAM kullanım bilgisini okur
	if err != nil {
		return CheckResult{}, fmt.Errorf("ram kullanimi okunamadi: %w", err)
	}

	percent := vmem.UsedPercent //kullanım yüzdesi okunur

	level := DetermineLevel(percent, r.ThresholdWarn, r.ThresholdError)

	var message string
	switch level {
	case LevelError:
		message = fmt.Sprintf("RAM kullanimi %.1f%%, hata esigi %.1f%% asildi", percent, r.ThresholdError)
	case LevelWarning:
		message = fmt.Sprintf("RAM kullanimi %.1f%%, uyari esigi %.1f%% asildi", percent, r.ThresholdWarn)
	default:
		message = fmt.Sprintf("RAM kullanimi %.1f%%", percent)
	}

	return CheckResult{
		CheckType:   "ram",
		Level:       level,
		Message:     message,
		MetricName:  "ram_used_percent",
		MetricValue: percent,
		MetricUnit:  "percent",
	}, nil
}
