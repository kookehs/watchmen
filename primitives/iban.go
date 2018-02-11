package primitives

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"math/big"
	"strconv"
)

const IBANSize = 34

type IBAN [IBANSize]byte

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

func (iban *IBAN) ConvertToInteger(b []byte) *big.Int {
	var buffer bytes.Buffer

	for i := 0; i < IBANSize; i++ {
		buffer.WriteString(strconv.Itoa(int(b[i])))
	}

	integer := new(big.Int)
	integer.SetString(buffer.String(), 10)
	return integer
}

func (iban *IBAN) SetChecksum(b []byte) {
	if len(b) == 2 {
		iban[2] = b[0]
		iban[3] = b[1]
	}
}

func (iban *IBAN) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(iban); err != nil {
		return err
	}

	return nil
}

func (iban *IBAN) DeseralizeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(iban); err != nil {
		return err
	}

	return nil
}

func (iban *IBAN) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(iban); err != nil {
		return err
	}

	return nil
}

func (iban *IBAN) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(iban); err != nil {
		return err
	}

	return nil
}

func (iban *IBAN) String() string {
	return string(iban[:])
}
