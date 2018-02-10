package primitives

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"io"
)

// Value of AccountHash is a sha256 hash
type AccountHash [sha256.Size]byte

type Account struct {
	Bban *BBAN `json:"bban"`
	Iban *IBAN `json:"iban"`
	Key  *Key  `josn:"key"`
}

func NewAccount(r io.Reader) (*Account, error) {
	key, err := NewKeyForICAP(r)

	if err != nil {
		return nil, err
	}

	bban := NewBBAN([]byte(key.Address.String()))
	iban := NewIBAN([]byte("TV00" + bban.String()))

	return &Account{
		Bban: bban,
		Iban: iban,
		Key:  key,
	}, nil
}

func (a *Account) Hash() [sha256.Size]byte {
	var buffer bytes.Buffer
	a.Serialize(&buffer)
	return sha256.Sum256(buffer.Bytes())
}

func (a *Account) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(a); err != nil {
		return err
	}

	return nil
}

func (a *Account) DeseralizeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(a); err != nil {
		return err
	}

	return nil
}

func (a *Account) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(a); err != nil {
		return err
	}

	return nil
}

func (a *Account) SeralizeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(a); err != nil {
		return err
	}

	return nil
}
