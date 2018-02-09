package base36

import (
	"math/big"
	"strings"
)

var base36 = []byte{
	'0', '1', '2', '3', '4', '5',
	'6', '7', '8', '9', 'A', 'B',
	'C', 'D', 'E', 'F', 'G', 'H',
	'I', 'J', 'K', 'L', 'M', 'N',
	'O', 'P', 'Q', 'R', 'S', 'T',
	'U', 'V', 'W', 'X', 'Y', 'Z',
}

func Encode(b []byte) []byte {
	if b[0] == '0' && b[1] == 'x' {
		b = b[2:]
	}

	integer := new(big.Int)
	integer.SetString(string(b), 16)
	return []byte(strings.ToUpper(integer.Text(36)))
}
