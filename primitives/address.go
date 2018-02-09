package primitives

import (
	"encoding/hex"
)

const AddressSize = 20

// Value of Address is a sha256 hash with the last 20 bytes
type Address [AddressSize]byte

func NewAddress(b []byte) *Address {
	address := new(Address)
	address.SetBytes(b)
	return address
}

func (a *Address) Hex() []byte {
	encoded := make([]byte, hex.EncodedLen(AddressSize))
	hex.Encode(encoded, (*a)[:])
	return encoded
}

func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressSize:]
	}

	copy(a[AddressSize-len(b):], b)
}

func (a *Address) String() string {
	return ("0x" + string(a.Hex()))
}
