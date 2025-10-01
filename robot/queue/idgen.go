package queue

import (
	"crypto/sha256"
)

func idgen(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	return h.Sum(nil)
}
