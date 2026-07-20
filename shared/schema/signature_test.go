package schema

import "testing"

func TestSignAndVerify_Success(t *testing.T) {
	secret := []byte("test-secret")
	b := Batch{
		BatchID:      "batch_1",
		HostID:       "web-03",
		AgentVersion: "0.1.0",
		CreatedAt:    "2026-07-16T09:14:05Z",
		Events: []Event{
			{EventID: "evt_1", HostID: "web-03", CheckType: "disk", Level: LevelError, Message: "disk full", Timestamp: "2026-07-16T09:14:03Z", SchemaVersion: 1},
		},
	}

	sig, err := Sign(secret, b)
	if err != nil {
		t.Fatalf("Sign hata verdi: %v", err)
	}
	b.Signature = sig

	if err := Verify(secret, b); err != nil {
		t.Errorf("Verify basarisiz olmamaliydi: %v", err)
	}
}

func TestVerify_WrongSecret_Fails(t *testing.T) {
	b := Batch{BatchID: "batch_1", HostID: "web-03"}

	sig, _ := Sign([]byte("dogru-anahtar"), b)
	b.Signature = sig

	err := Verify([]byte("yanlis-anahtar"), b)
	if err == nil {
		t.Error("yanlis anahtarla Verify basarili olmamaliydi, ama oldu")
	}
}

func TestVerify_TamperedData_Fails(t *testing.T) {
	secret := []byte("test-secret")
	b := Batch{BatchID: "batch_1", HostID: "web-03"}

	sig, _ := Sign(secret, b)
	b.Signature = sig

	b.HostID = "web-99"

	err := Verify(secret, b)
	if err == nil {
		t.Error("veri degistirildikten sonra Verify basarili olmamaliydi, ama oldu")
	}
}
