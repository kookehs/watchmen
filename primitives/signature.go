package primitives

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"math/big"
)

// Signature a pair of integers that represent the signature.
type Signature struct {
	R *big.Int
	S *big.Int
}

// MakeSignature creates and initializes a Signature from the given arguments.
func MakeSignature(r, s *big.Int) Signature {
	return Signature{
		R: r,
		S: s,
	}
}

// Deserialize decodes byte data encoded by gob.
func (s *Signature) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(s)
}

// DeseralizeJSON decodes JSON data.
func (s *Signature) DeseralizeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(s)
}

// Serialize encodes to byte data using gob.
func (s *Signature) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(s)
}

// SerializeJSON encodes to JSON data.
func (s *Signature) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(s)
}
