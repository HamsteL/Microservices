package middlewares

import (
	"crypto/sha256"
	"encoding/base64"
)

func GetStringHash(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
