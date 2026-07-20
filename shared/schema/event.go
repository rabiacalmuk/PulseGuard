package schema

type Level string

const (
	LevelInfo    Level = "INFO"
	LevelWarning Level = "WARNING"
	LevelError   Level = "ERROR"
)

type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"` // Dashboard'da grafik eksenini etiketlerken bu bilgi kullanılacak.
}

type Event struct {
	EventID       string `json:"event_id"`         // Olayın benzersiz kimliği
	HostID        string `json:"host_id"`          // olayın gerçekleştiği host'un benzersiz kimliği
	CheckType     string `json:"check_type"`       // olayın hangi türde kontrolile ilişkili olduğunu belirtir(disk,ram,cpu)
	Level         Level  `json:"level"`            // olayın önem derecesini belirtir, yukarıda tanımlı level tipini kullanır
	Message       string `json:"message"`          // olayın açıklaması, loglarda ve dashboard'da gösterilecek metin budur
	Metric        Metric `json:"metric,omitempty"` //olayla ilişkili sayısal ölçüm
	Timestamp     string `json:"timestamp"`        // olayın gerçekleştiği zaman, RFC3339 formatında
	SchemaVersion int    `json:"schema_version"`   // olayın hangi şema versiyonuyla üretildiğini belirtir,uyumluluk için kullanılır
}

//omiempty , metric alanına değer atanmadığı durumlarda json çıktısında metric alanının gözükmemesini sağlar.
