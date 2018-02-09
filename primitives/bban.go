package primitives

import (
	"github.com/kookehs/watchmen/encoding/base36"
)

const BBANSize = 30

type BBAN [BBANSize]byte

func NewBBAN(b []byte) *BBAN {
	bban := new(BBAN)
	bban.SetBytes(BBANFromHex(b))
	return bban
}

func (bban *BBAN) SetBytes(b []byte) {
	if len(b) > len(bban) {
		b = b[len(b)-BBANSize:]
	}

	copy(bban[BBANSize-len(b):], b)
}

func (bban *BBAN) String() string {
	return string(bban[:])
}

func BBANFromHex(b []byte) []byte {
	if (len(b) == 42) && (b[0] == '0') && (b[1] == 'x') {
		b[2] = '0'
		b[3] = 'x'
		b = b[2:]
		return base36.Encode(b)
	}

	return nil
}
