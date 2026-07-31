package sender

import (
	"time"

	"pulseguard/agent/internal/checks"
	"pulseguard/shared/schema"
)

func ToEvent(hostID string, r checks.CheckResult) schema.Event { //checks paketinin üretimini shema paketinin beklediği formata çeviren fonksiyon
	return schema.Event{
		EventID:   generateID(),
		HostID:    hostID,
		CheckType: r.CheckType,
		Level:     schema.Level(r.Level), //checks.level ve schema.level aynı isimdde farklı tiplerdi, burada birini diğerine dönüştürerek tip dönüşümü yapıyoruz
		Message:   r.Message,
		Metric: schema.Metric{
			Name:  r.MetricName,
			Value: r.MetricValue,
			Unit:  r.MetricUnit,
		},
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		SchemaVersion: 1,
	}
}
