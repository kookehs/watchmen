package watchmen

import (
	"math/big"
)

type Amount big.Int

// Hash values of accounts are 256 bits
type AccountHash [32]byte
