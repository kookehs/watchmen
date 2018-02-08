package primitives

import (
	"crypto/sha256"
	"math/big"
)

// Value of Address is a sha256 hash
type Address [sha256.Size]byte

type Amount = big.Int

// Value of AccountHash is a sha256 hash
type AccountHash [sha256.Size]byte

type Account struct {
	Key Key
}
