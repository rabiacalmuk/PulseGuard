package checks

import (
	"fmt" //hata mesajlarını ve string formatlamayı yapmak için bu paketi kullanıyoruz

	"github.com/shirou/gopsutil/v3/disk" //disk kullanımını okumak için bu paketi kullanıyoruz
)

type DiskCheck struct { //check'in ihtiyaç duyduğu ayarları taşıyan struct
	Mount          string  //kontrol edilecek disk bölümü
	ThresholdWarn  float64 //yüzde kaçtan itibaren WARNING seviyesine geçileceğini belirtir
	ThresholdError float64 //yüzde kaçtan itibaren ERROR seviyesine geçileceğini belirtir
}

func (d DiskCheck) Name() string { // check tipinin adını döndüren metot
	return "disk" //check tipinin adı "disk" olarak belirlenmiştir
}

func (d DiskCheck) Run() (CheckResult, error) { //check tipinin tek bir ölçümden ürettiği ham sonucu döndüren metot
	usage, err := disk.Usage(d.Mount) //gopsutil'in disk.Usage fonksiyonu çağrılır, bu OS'den gerçek disk kullanım bilgisini okur
	if err != nil {
		return CheckResult{}, fmt.Errorf("disk usage okunamadi: %w", err)
	} //mount noktası geçersizse veya disk bilgisi okunamazsa boş bir CheckResult ve bir hata mesajı döndürülür, orijinal hatayı sarmalayarak saklarız

	percent := usage.UsedPercent //kullanım yüzdesi okunur

	level := DetermineLevel(percent, d.ThresholdWarn, d.ThresholdError)

	var message string
	switch level {
	case LevelError:
		message = fmt.Sprintf("Disk kullanimi %.1f%%, hata esigi %.1f%% asildi (%s)", percent, d.ThresholdError, d.Mount)
	case LevelWarning:
		message = fmt.Sprintf("Disk kullanimi %.1f%%, uyari esigi %.1f%% asildi (%s)", percent, d.ThresholdWarn, d.Mount)
	default:
		message = fmt.Sprintf("Disk kullanimi %.1f%% (%s)", percent, d.Mount)
	}

	return CheckResult{
		CheckType:   "disk",
		Level:       level,
		Message:     message,
		MetricName:  "disk_used_percent",
		MetricValue: percent,
		MetricUnit:  "percent",
	}, nil //CheckResult döndürülüyor, hata yoksa nil döndürülüyor
}
