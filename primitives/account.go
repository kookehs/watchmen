package watchmen

import (
	"crypto/sha256"
	"math/big"
)

type Amount = big.Int

// Hash values of accounts are 256 bits
type AccountHash [sha256.Size]byte
