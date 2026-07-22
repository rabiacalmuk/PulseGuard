package checks

func DetermineLevel(percent, thresholdWarn, thresholdError float64) Level {
	switch {
	case percent >= thresholdError:
		return LevelError
	case percent >= thresholdWarn:
		return LevelWarning
	default:
		return LevelInfo
	}
}

//disk,cpu ve ram  dosyalarının üçünde de tekrar eden switch case yapısını burada topluyoruz test edebilmek için
