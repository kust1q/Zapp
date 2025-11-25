package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Hasher struct {
	secret []byte
}

func NewHasher(secret string) *Hasher {
	return &Hasher{
		secret: []byte(secret),
	}
}

func (h *Hasher) AuthHash(dataType, data string) string {
	mac := hmac.New(sha256.New, h.secret)
	mac.Write([]byte(dataType + data))
	return hex.EncodeToString(mac.Sum(nil))
}
