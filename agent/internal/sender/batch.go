package sender

import (
	"time"

	"pulseguard/shared/schema"
)

func BuildBatch(secret []byte, hostID, agentVersion string, events []schema.Event) (schema.Batch, error) {
	b := schema.Batch{
		BatchID:      generateID(),
		HostID:       hostID,
		AgentVersion: agentVersion,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
		Events:       events,
	}

	sig, err := schema.Sign(secret, b) //daha önce signature.go'da yazdığımız sign fonksiyonunu çağırıyoruz o da batch'i imzalıyor
	if err != nil {
		return schema.Batch{}, err
	}
	b.Signature = sig //hesaplanan imza batch'in signature alanına yerleştirilir ve imzalı bir batch'e sahip olmuş oluruz artık

	return b, nil
}
