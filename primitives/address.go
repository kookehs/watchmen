package primitives

import (
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"io"
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

func (a *Address) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(a); err != nil {
		return err
	}

	return nil
}

func (a *Address) DeseralizeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(a); err != nil {
		return err
	}

	return nil
}

func (a *Address) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(a); err != nil {
		return err
	}

	return nil
}

func (a *Address) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(a); err != nil {
		return err
	}

	return nil
}

func (a *Address) String() string {
	return ("0x" + string(a.Hex()))
}
