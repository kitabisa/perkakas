package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateHmac will transform data into HMAC signature based on sha256 & secret_key
func GenerateHmac(data, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))

	result := h.Sum(nil)
	return hex.EncodeToString(result)
}

// IsMatchHmac will check whether signature is match with expected HMAC signature
func IsMatchHmac(data, signature, secretKey string) bool {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))

	expected := h.Sum(nil)
	sign, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if hmac.Equal(expected, sign) {
		return true
	}
	return false
}
