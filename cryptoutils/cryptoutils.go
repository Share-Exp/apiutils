package cryptoutils

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMD5 is a function to cryptograph a string.
func GetMD5(input string) string {
	hash := md5.New()
	hash.Reset()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}
