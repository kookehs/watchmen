package primitives

import (
	"encoding/gob"
	"encoding/json"
	"io"

	"github.com/kookehs/watchmen/encoding/base36"
)

// BBANSize is a the fixed length of a Basic Bank Account Number
const BBANSize = 30

// BBAN contains up to 30 alphanumeric characters
type BBAN [BBANSize]byte

// MakeBBAN creates and initializes a BBAN from the given bytes.
func MakeBBAN(b []byte) BBAN {
	var bban BBAN
	bban.SetBytes(BBANFromHex(b))
	return bban
}

// SetBytes sets the bytes of the BBAN to the given bytes.
func (bban *BBAN) SetBytes(b []byte) {
	if len(b) > len(bban) {
		b = b[len(b)-BBANSize:]
	}

	copy(bban[BBANSize-len(b):], b)
}

// Deserialize decodes byte data encoded by gob.
func (bban *BBAN) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(bban)
}

// DeseralizeJSON decodes JSON data.
func (bban *BBAN) DeseralizeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(bban)
}

// Serialize encodes to byte data using gob.
func (bban *BBAN) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(bban)
}

// SerializeJSON encodes to JSON data.
func (bban *BBAN) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(bban)
}

// String returns the string representation of the BBAN.
func (bban *BBAN) String() string {
	return string(bban[:])
}

// BBANFromHex generates a BBAN from the given bytes representing a hex number.
func BBANFromHex(b []byte) []byte {
	if (len(b) == 42) && (b[0] == '0') && (b[1] == 'x') {
		b[2] = '0'
		b[3] = 'x'
		b = b[2:]
		encoded := base36.Encode(b)

		if padding := IBANSize - len(encoded); padding > 0 {
			bban := make([]byte, IBANSize)

			for i := 0; i < padding; i++ {
				bban = append(bban, '0')
			}

			bban = append(bban, encoded...)
			return bban
		}

	}

	return nil
}
