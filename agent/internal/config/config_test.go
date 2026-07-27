package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ValidConfig(t *testing.T) { //geçerli bir YAML verildiğinde load doğru config struct'ını veriyor mu bu fonksiyonla kontrol ederiz
	dir := t.TempDir()                              //teste özel, test sonunda otomatik silinecek geçici bir klasör oluşturuyor
	configPath := filepath.Join(dir, "config.yaml") // tam dosya yolu belirlenir

	yamlContent := `
collector:
  url: "https://collector.internal:8443"
  auth_token: "sabit-token-123"
  batch_interval_seconds: 30
  retry_backoff_seconds: [5, 15, 60]

host:
  id: "web-03"
  tags: ["prod", "web-tier"]

checks:
  - type: disk
    interval_seconds: 60
    mount: "/"
    threshold_warning: 85
    threshold_error: 92

queue:
  path: "/var/lib/pulseguard/queue.db"
  max_size_mb: 200
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("gecici config dosyasi yazilamadi: %v", err)
	} //yamlContent, configPath'e gerçekten yazılır. 0644 dosya izinlerinde sahibin hem okuma hem yazma yapabileceğini diğerlerininse sadece okuma yapabileceğini belirten bir varsayılandır.

	cfg, err := Load(configPath) //gerçek load fonk gerçek bir dosya yoluyla çağırılır
	if err != nil {
		t.Fatalf("Load hata verdi: %v", err)
	}

	if cfg.Collector.URL != "https://collector.internal:8443" {
		t.Errorf("Collector.URL yanlis: %s", cfg.Collector.URL)
	} //yaml url alanının doğru şekilde yerleştiği kontrol edilir
	if cfg.Host.ID != "web-03" {
		t.Errorf("Host.ID yanlis: %s", cfg.Host.ID)
	}
	if len(cfg.Checks) != 1 {
		t.Fatalf("Checks sayisi yanlis: %d", len(cfg.Checks))
	} //yaml'da tek bir check tanımlamıştık bu yüzden liste uzunluğu 1 olmalı.Fatalf kullanmamızın nedeni de eğer sayı yanlışsa bu cfg.Checks[0] erişiminde panic'e neden olabilir bunu önlemek için de fatlf kullandık.
	if cfg.Checks[0].Type != "disk" {
		t.Errorf("Checks[0].Type yanlis: %s", cfg.Checks[0].Type)
	}
	if cfg.Checks[0].ThresholdError != 92 {
		t.Errorf("Checks[0].ThresholdError yanlis: %v", cfg.Checks[0].ThresholdError)
	} //listenin içindeki check'in alanlarının doğru yerleştiğğini kontrol ediyoruz
}

func TestLoad_FileNotFound(t *testing.T) { //hata senaryosu, olmayan bir dosya verildiğinde load gerçekten hata döndürüyor mu kontrol edilir
	_, err := Load("olmayan-dosya.yaml") //_ config sonucunu önemsemediğimizi belirtir, sadece err sonucuna bakıyoruz
	if err == nil {
		t.Error("olmayan dosya icin Load hata dondurmeliydi, ama donmedi")
	} // err== nil ise olmayan dosyayı başarıyla okumuş gibi davranıyor ve bu bir problem demek olur
}

func TestResolveEnv_WithEnvVar(t *testing.T) { //resolveEnv fonksiyonu ${...} kalıbını doğru çözüyor mu diye bakılır
	t.Setenv("PULSEGUARD_TEST_TOKEN", "gizli-deger-123") //test süresince bu env değişkeni geçerli olur test sonunda go otomatik geri alır bu değişkeni

	result := resolveEnv("${PULSEGUARD_TEST_TOKEN}") //${} kaldırıp setenv ile tanımladığımız değeri döndürecek

	if result != "gizli-deger-123" {
		t.Errorf("beklenen 'gizli-deger-123', gelen: %s", result)
	}
}

func TestResolveEnv_PlainValue(t *testing.T) { //resolveEnv'in ${} kalıbında olmayan düz bir değeri olduğu gibi bıraktığı test edilir
	result := resolveEnv("duz-bir-deger") //değer değiştirilmeden geri dönmeli

	if result != "duz-bir-deger" {
		t.Errorf("kalip disindaki deger degismemeliydi, gelen: %s", result)
	}
}
