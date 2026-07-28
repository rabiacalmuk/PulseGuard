package queue

import (
	"path/filepath"
	"testing"

	"pulseguard/shared/schema"
)

func TestEnqueueAndAll(t *testing.T) { //iki event ekleyip All() ile okuma yapıyor ve ikisi de doğru şekilde geri geliyor mu kontrol ediyoruz
	path := filepath.Join(t.TempDir(), "queue.db") //test sonunda otomatik olarak silinecek bir klasör içine dosya yolu oluşturuyoruz
	q, err := Open(path)
	if err != nil {
		t.Fatalf("Open hata verdi: %v", err)
	}

	e1 := schema.Event{EventID: "evt_1", HostID: "web-03", CheckType: "disk", Level: schema.LevelError, Message: "disk full", Timestamp: "2026-07-16T09:14:03Z", SchemaVersion: 1}
	e2 := schema.Event{EventID: "evt_2", HostID: "web-03", CheckType: "cpu", Level: schema.LevelWarning, Message: "cpu high", Timestamp: "2026-07-16T09:15:00Z", SchemaVersion: 1}

	if err := q.Enqueue(e1); err != nil {
		t.Fatalf("Enqueue hata verdi: %v", err)
	}
	if err := q.Enqueue(e2); err != nil {
		t.Fatalf("Enqueue hata verdi: %v", err)
	}

	events, err := q.All()
	if err != nil {
		t.Fatalf("All hata verdi: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("event sayisi yanlis: %d", len(events))
	}
	if events[0].EventID != "evt_1" || events[1].EventID != "evt_2" {
		t.Errorf("event sirasi/icerigi yanlis: %+v", events)
		//ekleme sırasının korunduğu kontrol ediliyor, okuma sırası ekleme sırasıyla aynı olmalı
	}
}

func TestAck_RemovesOnlySpecifiedEvents(t *testing.T) { //ack çağırdığımızda sadece belirtilen event mi siliniyor diğrei kalıyor mu kontrol ediyoruz
	path := filepath.Join(t.TempDir(), "queue.db")
	q, _ := Open(path)

	e1 := schema.Event{EventID: "evt_1", HostID: "web-03"}
	e2 := schema.Event{EventID: "evt_2", HostID: "web-03"}
	q.Enqueue(e1)
	q.Enqueue(e2)

	if err := q.Ack([]string{"evt_1"}); err != nil {
		t.Fatalf("Ack hata verdi: %v", err)
	}

	events, err := q.All()
	if err != nil {
		t.Fatalf("All hata verdi: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("event sayisi yanlis: %d", len(events))
	}
	if events[0].EventID != "evt_2" {
		t.Errorf("yanlis event kaldi: %s", events[0].EventID) //kalan event'in doğru olan olduğunu kontrol ediyoruz
	}
}
