package config

import (
	"fmt"
	"os" // dosya ve env değişkeni okumak için
	"strings"

	"gopkg.in/yaml.v3"
)

type CollectorConfig struct {
	URL                  string `yaml:"url"`        //yaml:url tagi Json'daki gibi yaml alanındaki url anahtarının go'daki URL anahtarına bağlanmasını sağlar
	AuthToken            string `yaml:"auth_token"` //şablon olarak duracak load fonksiyonu ile bunu gerçek değer ile değiştireceğiz
	BatchIntervalSeconds int    `yaml:"batch_interval_seconds"`
	RetryBackoffSeconds  []int  `yaml:"retry_backoff_seconds"`
}

type HostConfig struct { //yamldaki host bloğu karşılığı
	ID   string   `yaml:"id"`
	Tags []string `yaml:"tags"`
}

type CheckConfig struct { //checks listesinin her bir elemanının karşılığı,yaml'da ne yazıldığı bilgisini taşıyan ham veri
	Type             string  `yaml:"type"`
	IntervalSeconds  int     `yaml:"interval_seconds"`
	ThresholdWarning float64 `yaml:"threshold_warning"`
	ThresholdError   float64 `yaml:"threshold_error"`
	Mount            string  `yaml:"mount"` //sadece disk type check için kullanılacak varsayılanıboş string, cpu ve ram tipi checkler oluşturulurken hiç bakılmayacak disk tipi checklerde okunacak
}

type QueueConfig struct { //yaml queue bloğu karşılığı
	Path      string `yaml:"path"`
	MaxSizeMB int    `yaml:"max_size_mb"`
}

type Config struct { //tüm altyapıların bir araya geldiği kısım, yaml dosyasının en temel 4 ana anahtarına karşılık geliyor
	Collector CollectorConfig `yaml:"collector"`
	Host      HostConfig      `yaml:"host"`
	Checks    []CheckConfig   `yaml:"checks"`
	Queue     QueueConfig     `yaml:"queue"`
}

func Load(path string) (Config, error) { //dışarıdan çağırılacak ana fonksiyon,yaml dosyasının nerede olduğunu belirtir
	data, err := os.ReadFile(path) //dosya ham byte dizisi olarak okunuyor , henüz yorumlanmamış halde sadece dosya içeriği
	if err != nil {
		return Config{}, fmt.Errorf("config dosyasi okunamadi (%s): %w", path, err)
	} //dosya açılamıyor veya okunamıyorsa açıklayıcı bir hata döndürülüyor ve hangi dosyanın sorunlu olduğunun bilinmesi açısından path bilgisi de mesaja ekleniyor

	var cfg Config                                     //yaml verisini buraya doldurmak için boş bir config struct oluşturuyoruz
	if err := yaml.Unmarshal(data, &cfg); err != nil { //yaml.Marshal ile ham byte dizisini alıp struct taglerine bakarak doğru alanlara yerleştirme işlemini yaptırıyoruz
		return Config{}, fmt.Errorf("config parse edilemedi: %w", err) //yaml kendisi bozuksa bu hata döndürülür
	}

	cfg.Collector.AuthToken = resolveEnv(cfg.Collector.AuthToken) //yaml'dan "${PULSEGUARD_AGENT_TOKEN}" gibi şablon bir değer geldiyse bunu gerçek env değişkeniyle değiştiriyoruz

	return cfg, nil //başarıyla doldurulmuş config döner
}

func resolveEnv(value string) string { //yaml'den egelen bir stringin "${ENV_ADI}" kalıbında olup olmadığı kontrolünü yapan fonksiyon
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") { //"${" ile başlayıp "{" ile bitiyor mu kontrol ediyoruz
		envName := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}") //"${" ve "{" kısımları kırpılıyor ifadeden
		return os.Getenv(envName)                                           //kırptığımız halde kalan ifade ile isimlendirilmiş bir env değişkeni var mı diye bakılıyor varsa okuyup döndürülüyor yoksa boş string döndürülüyor
	}
	return value //"${}" kalıbında değilse değer olduğu gibi dönüyor
}
