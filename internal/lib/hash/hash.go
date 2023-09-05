package hash

import (
	"crypto/sha1"
	"encoding/base64"
)

func Hash(pass string) string {
	hasher := sha1.New()
	hasher.Write([]byte(pass))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
