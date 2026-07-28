package sender

import (
	"testing"

	"pulseguard/shared/schema"
)

func TestBuildBatch_ProducesValidSignature(t *testing.T) { //buildBatch'in ürettiği batch'in hem doğru alanlara sahip olduğunu hem de geçerli bir imzası olduğunun kontrolünü yapıyoruz burada
	secret := []byte("test-secret")
	events := []schema.Event{
		{EventID: "evt_1", HostID: "web-03", CheckType: "disk", Level: schema.LevelError, Message: "disk full", Timestamp: "2026-07-16T09:14:03Z", SchemaVersion: 1},
	}

	b, err := BuildBatch(secret, "web-03", "0.1.0", events)
	if err != nil {
		t.Fatalf("BuildBatch hata verdi: %v", err)
	}

	if b.BatchID == "" {
		t.Error("BatchID bos olmamaliydi")
	}
	if b.HostID != "web-03" {
		t.Errorf("HostID yanlis: %s", b.HostID)
	}
	if b.AgentVersion != "0.1.0" {
		t.Errorf("AgentVersion yanlis: %s", b.AgentVersion)
	}
	if len(b.Events) != 1 {
		t.Fatalf("Events sayisi yanlis: %d", len(b.Events))
	}

	if err := schema.Verify(secret, b); err != nil {
		t.Errorf("Verify basarisiz olmamaliydi: %v", err)
	}
}

func TestBuildBatch_UniqueIDs(t *testing.T) { //iki farklı Buildbatch çağrısı yapılır gerçekten farklı batchID'ler üretiyor mu diye kontrol etmek için
	secret := []byte("test-secret")

	b1, _ := BuildBatch(secret, "web-03", "0.1.0", nil)
	b2, _ := BuildBatch(secret, "web-03", "0.1.0", nil)

	if b1.BatchID == b2.BatchID {
		t.Error("iki farkli BuildBatch cagrisi ayni BatchID'yi uretti")
	}
}
