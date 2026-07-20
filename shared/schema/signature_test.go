package schema //test dosyaları test ettikleri kodla aynı pakette olmalı

import "testing"

func TestSignAndVerify_Success(t *testing.T) { //sign ve verify fonksiyonlarının doğru çalıştığını test eder
	secret := []byte("test-secret") //test edilecek dosyalar byte dizisi bekliyordu bu yüzden stringi byte'a çeviriyoruz
	b := Batch{                     //test için sahte batch oluşturuluyor
		BatchID:      "batch_1",
		HostID:       "web-03",
		AgentVersion: "0.1.0",
		CreatedAt:    "2026-07-16T09:14:05Z",
		Events: []Event{ //batch en az bir event içermeli bu yüzden bir örnek event eklenir
			{EventID: "evt_1", HostID: "web-03", CheckType: "disk", Level: LevelError, Message: "disk full", Timestamp: "2026-07-16T09:14:03Z", SchemaVersion: 1},
		},
	}

	sig, err := Sign(secret, b) //batch'i imzalar ve imza değerini döndürür
	if err != nil {
		t.Fatalf("Sign hata verdi: %v", err) //t.Fatalf, testin başarısız olduğunu bildirir ve testin geri kalanını çalıştırmaz, sign patladıktan sonra verify çalıştırmak anlamsız. Kritik durumlar için t.Fatalf kullanılır.
	}
	b.Signature = sig //batch'in Signature alanına imza değerini ekler

	if err := Verify(secret, b); err != nil { //batch'i doğrular, eğer imza geçersizse hata döndürür
		t.Errorf("Verify basarisiz olmamaliydi: %v", err) //t.Errorf testi başarısız işaretler ama durdurmaz , kritik olmayan kontroller için kullanlır.
	}
}

func TestVerify_WrongSecret_Fails(t *testing.T) { //yanlış anahtarla doğrulama denenirse ne olacağını test ederiz. Başarılı olmamalı.
	b := Batch{BatchID: "batch_1", HostID: "web-03"}

	sig, _ := Sign([]byte("dogru-anahtar"), b)
	b.Signature = sig //doğru anahtarla üretilen imza batch'e yerleştiriliyor

	err := Verify([]byte("yanlis-anahtar"), b) //yanlış anahtarla doğrulama deneniyor, err nil olmamalı
	if err == nil {                            //yanlış anahtarla doğrulama başarılı olursa test başarısız olur
		t.Error("yanlis anahtarla Verify basarili olmamaliydi, ama oldu")
	}
}

func TestVerify_TamperedData_Fails(t *testing.T) { //veri imzalandıktan sonra ne olacağı test edilir, verify bunu yakalamalı
	secret := []byte("test-secret")
	b := Batch{BatchID: "batch_1", HostID: "web-03"}

	sig, _ := Sign(secret, b) //batch henüz değiştirilmemiş haliyle imzalanıyor
	b.Signature = sig

	b.HostID = "web-99" //batch'in HostID alanı değiştirilerek veri bozuluyor, imza artık geçersiz olmalı

	err := Verify(secret, b) //veri değiştirildikten sonra doğrulama deneniyor, err nil olmamalı
	if err == nil {          //veri değiştirildikten sonra doğrulama başarılı olursa test başarısız olur
		t.Error("veri degistirildikten sonra Verify basarili olmamaliydi, ama oldu")
	}
}
