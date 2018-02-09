package primitives

import (
	"crypto/sha256"
	"io"
)

// Value of AccountHash is a sha256 hash
type AccountHash [sha256.Size]byte

type Account struct {
	Bban *BBAN
	Iban *IBAN
	Key  *Key
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
