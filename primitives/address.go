package primitives

import (
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"io"
)

// AddressSize is the fixed length of addresses
const AddressSize = 20

// Address a sha256 hash with the last 20 bytes
type Address [AddressSize]byte

// MakeAddress creates and initializes an Address structure with the given bytes.
func MakeAddress(b []byte) Address {
	var address Address
	address.SetBytes(b)
	return address
}

// Hex returns the hex representation of the Address.
func (a *Address) Hex() []byte {
	encoded := make([]byte, hex.EncodedLen(AddressSize))
	hex.Encode(encoded, (*a)[:])
	return encoded
}

// SetBytes sets the bytes of the Address to the given bytes.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressSize:]
	}

	copy(a[AddressSize-len(b):], b)
}

// Deserialize decodes byte data encoded by gob.
func (a *Address) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(a)
}

// DeseralizeJSON decodes JSON data.
func (a *Address) DeseralizeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(a)
}

// Serialize encodes to byte data using gob.
func (a *Address) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(a)
}

// SerializeJSON encodes to JSON data.
func (a *Address) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(a)
}

// String returns the hex representation of the Address with 0x prepended.
func (a *Address) String() string {
	return ("0x" + string(a.Hex()))
}
