package index

import (
	"crypto/sha256"
)

func IDGen(b []byte) [32]byte {
	var x = sha256.Sum256(b)
	return x
}
