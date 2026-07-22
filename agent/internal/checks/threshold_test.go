package checks

import "testing"

func TestDetermineLevel(t *testing.T) {
	tests := []struct {
		name           string //testin açıklaması
		percent        float64
		thresholdWarn  float64
		thresholdError float64
		want           Level //hangi level'in dönmesi bekleniyor
	}{
		{"esiklerin altinda", 50, 75, 90, LevelInfo},
		{"warning esigine tam esit", 75, 75, 90, LevelWarning},
		{"warning ile error arasinda", 80, 75, 90, LevelWarning},
		{"error esigine tam esit", 90, 75, 90, LevelError},
		{"error esiginin uzerinde", 95, 75, 90, LevelError},
		{"sifir yuzde", 0, 75, 90, LevelInfo},
	}

	for _, tt := range tests { //tablodaki tüm test caseler için test çalıştırılır
		t.Run(tt.name, func(t *testing.T) { //her test case için alt test çalıştırılır
			got := DetermineLevel(tt.percent, tt.thresholdWarn, tt.thresholdError)
			if got != tt.want {
				t.Errorf("DetermineLevel(%v, %v, %v) = %v, istenen %v",
					tt.percent, tt.thresholdWarn, tt.thresholdError, got, tt.want)
			} //DetermineLevel fonksiyonunun döndürdüğü level ile istenen level karşılaştırılır, eşleşmezse hata mesajı yazdırılır
		})
	}
}
