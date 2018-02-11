package primitives

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func MakeSignature(r, s *big.Int) Signature {
	return Signature{
		R: r,
		S: s,
	}
}

func (s *Signature) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(s); err != nil {
		return err
	}

	return nil
}

func (s *Signature) DeseralizeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(s); err != nil {
		return err
	}

	return nil
}

func (s *Signature) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(s); err != nil {
		return err
	}

	return nil
}

func (s *Signature) SeralizeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(s); err != nil {
		return err
	}

	return nil
}
