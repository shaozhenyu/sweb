package apikey

import (
	"crypto/rand"
)

func Gen(apikeyLen int) []byte {
	b := make([]byte, apikeyLen)
	rand.Read(b)
	return b
}
