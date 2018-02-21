package primitives

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"math/big"
	"strconv"
)

// IBANSize is the fixed length of an International Bank Account Number
const IBANSize = 34

// IBAN represents an International Bank Account Number
// which consists of up to 34 alphanumeric characters.
// Country Code - 2 bytes
// Checksum - 2 bytes
// BBAN - 30 bytes
type IBAN [IBANSize]byte

// MakeIBAN creates and initializes an IBAN from the given bytes.
// Takes a slice of bytes in the format of an IBAN.
func MakeIBAN(b []byte) IBAN {
	var iban IBAN
	copy(iban[:], b)
	copy(b[IBANSize-4:], iban[:4])
	copy(b[:IBANSize-4], iban[4:])
	numeric := iban.ConvertToNumeric(b)
	integer := iban.ConvertToInteger(numeric)
	checksum := iban.CalculateChecksum(integer)
	iban.SetChecksum(checksum)
	return iban
}

// CalculateChecksum calcuates the checksum of the IBAN.
func (iban *IBAN) CalculateChecksum(i *big.Int) []byte {
	mod := big.NewInt(97)
	remainder := new(big.Int)
	remainder.Mod(i, mod)
	checksum := 98 - int(remainder.Int64())
	var buffer bytes.Buffer

	if checksum < 10 {
		buffer.WriteByte('0')
	}

	buffer.WriteString(strconv.Itoa(checksum))
	return buffer.Bytes()
}

// ConvertToNumeric converts all letters to their numeric values.
func (iban *IBAN) ConvertToNumeric(b []byte) []byte {
	for i := IBANSize - 1; i >= 0; i-- {
		if (b[i] >= 65) && (b[i] <= 90) {
			b[i] = b[i] - 65 + 10
		} else {
			b[i] = b[i] - 48
		}
	}

	return b
}

// ConvertToInteger takes the numeric representation in bytes and converts it to a big.Int.
func (iban *IBAN) ConvertToInteger(b []byte) *big.Int {
	var buffer bytes.Buffer

	for i := 0; i < IBANSize; i++ {
		buffer.WriteString(strconv.Itoa(int(b[i])))
	}

	integer := new(big.Int)
	integer.SetString(buffer.String(), 10)
	return integer
}

// SetChecksum sets the 2 bytes allocated for a checksum to the given bytes.
func (iban *IBAN) SetChecksum(b []byte) {
	if len(b) == 2 {
		iban[2] = b[0]
		iban[3] = b[1]
	}
}

// Deserialize decodes byte data encoded by gob.
func (iban *IBAN) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(iban)
}

// DeseralizeJSON decodes JSON data.
func (iban *IBAN) DeseralizeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(iban)
}

// Serialize encodes to byte data using gob.
func (iban *IBAN) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(iban)
}

// SerializeJSON encodes to JSON data.
func (iban *IBAN) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(iban)
}

// String returns the string representation of the IBAN.
func (iban *IBAN) String() string {
	return string(iban[:])
}
