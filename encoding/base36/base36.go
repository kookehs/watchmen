package base36

import (
	"math/big"
	"strings"
)

// Encode returns the base36 encoding of the given bytes.
func Encode(b []byte) []byte {
	if b[0] == '0' && b[1] == 'x' {
		b = b[2:]
	}

	integer := new(big.Int)
	integer.SetString(string(b), 16)
	return []byte(strings.ToUpper(integer.Text(36)))
}
