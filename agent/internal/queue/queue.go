package queue

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"pulseguard/shared/schema"
)

type Queue struct {
	path string
	mu   sync.Mutex //aynı anda sadece bir goroutine'in kuyruğa erişebilmesini sağlar
}

func Open(path string) (*Queue, error) { //kuyruğu açan fonksiyondur, dosya yoksa oluşturur varsa içeriğini değiştirmez
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("kuyruk dosyasi acilamadi (%s): %w", path, err)
	}
	f.Close()
	return &Queue{path: path}, nil //çağıran tarafın aynı queue nesnesi üzerinde işlem yapabilmesi için bir pointer döndürüyoruz
}

func (q *Queue) Enqueue(event schema.Event) error { //kuyruğa yeni bir event ekleyen metottur
	q.mu.Lock()         //kilidi alır, bu satırdan fonksiyon bitene kadarhiçbir goroutine aynı kilide erişemez
	defer q.mu.Unlock() //defer ile fonksiyon bittiğinde (sonucu önemsiz) otomatik olarak kilidi bırakır

	f, err := os.OpenFile(q.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //dosya sonuna ekleme modunda açılır , var olan içerik silinmeden yeni satır en sona eklenir
	if err != nil {
		return fmt.Errorf("kuyruk dosyasi acilamadi: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(event) //event struct'ı json byte dizisine çevrilir
	if err != nil {
		return fmt.Errorf("event serialize edilemedi: %w", err)
	}

	if _, err := f.Write(append(data, '\n')); err != nil { //json verisinin sonuna bir yeni satır karakteri eklenir böylece her satır tek bir event'i temsil eder
		return fmt.Errorf("event yazilamadi: %w", err)
	}
	return nil
}

func (q *Queue) All() ([]schema.Event, error) { //kuyruktaki tüm eventleri okuyan sışarıya açık metod
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.readAll() //gerçek okuma işlemi readAll fonksiyonuna devredildi
}

func (q *Queue) Ack(eventIDs []string) error { //başarıyla gönderilmiş eventleri kuyruktan silen metod, collector'a başarıyla ulaşan eventlerin ID'lerini alır
	q.mu.Lock()
	defer q.mu.Unlock()

	events, err := q.readAll()
	if err != nil {
		return err
	}

	acked := make(map[string]bool) //gönderilen ID'ler map haline getirilir
	for _, id := range eventIDs {
		acked[id] = true
	}

	var remaining []schema.Event //ack'lenmemiş eventleri koruyup diğerlerini atıyoruz
	for _, e := range events {
		if !acked[e.EventID] {
			remaining = append(remaining, e)
		}
	}

	return q.rewrite(remaining) //dosyayı kalan henüz gönderilmemiş olan eventlerle yeniden yazıyoruz
}

func (q *Queue) readAll() ([]schema.Event, error) { //kilit almayan (private) yardımcı fonksiyon, dosyayı ssatır satır okuyup her satırı bir event'e çevirir
	f, err := os.Open(q.path)
	if err != nil {
		return nil, fmt.Errorf("kuyruk dosyasi acilamadi: %w", err)
	}
	defer f.Close()

	var events []schema.Event
	scanner := bufio.NewScanner(f) //tüm dosyayı tek seferde belleğe yüklemek yerine akış halinde okumamızı sağlar scanner
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var e schema.Event
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("event parse edilemedi: %w", err)
		}
		events = append(events, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("kuyruk dosyasi okunamadi: %w", err)
	}
	return events, nil
}

func (q *Queue) rewrite(events []schema.Event) error { // private yardımcı fonksiyon dosyanın tamamını silip verilen event listesini yeniden oluşturur
	f, err := os.Create(q.path) //dosya varsa içeriğini sıfırlar yoksa oluşturur
	if err != nil {
		return fmt.Errorf("kuyruk dosyasi yeniden yazilamadi: %w", err)
	}
	defer f.Close()

	for _, e := range events {
		data, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("event serialize edilemedi: %w", err)
		}
		if _, err := f.Write(append(data, '\n')); err != nil {
			return fmt.Errorf("event yazilamadi: %w", err)
		}
	}
	return nil
}
