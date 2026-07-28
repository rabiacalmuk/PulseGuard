package sender

import (
	"crypto/rand"
	"encoding/hex"
)

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

//id'ler için rastgele verilen uzunlukta rastgele veri üretiyor ve bu ham veriyi okunabilir hex stringe çeviriyoruz
