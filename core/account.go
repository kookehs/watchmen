package core

import (
	"encoding/gob"
	"encoding/json"
	"io"

	"github.com/kookehs/watchmen/primitives"
)

// Account contains address as well as the key that generated it.
type Account struct {
	Bban primitives.BBAN `json:"bban"`
	Iban primitives.IBAN `json:"iban"`
	Key  primitives.Key  `josn:"key"`
}

// MakeAccount creates and initializes an account with the given key.
func MakeAccount(key primitives.Key) Account {
	bban := primitives.MakeBBAN([]byte(key.Address.String()))
	iban := primitives.MakeIBAN([]byte("TV00" + bban.String()))

	return Account{
		Bban: bban,
		Iban: iban,
		Key:  key,
	}
}

// CreateChangeBlock creates a signed ChangeBlock with the given arguments.
func (a *Account) CreateChangeBlock(amt primitives.Amount, delegates []primitives.IBAN, prev primitives.BlockHash) (*primitives.ChangeBlock, error) {
	block := primitives.NewChangeBlock(amt, delegates, prev)

	if err := block.Sign(a.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

// CreateOpenBlock creates a signed OpenBlock.
func (a *Account) CreateOpenBlock() (*primitives.OpenBlock, error) {
	var balance primitives.Amount
	block := primitives.NewOpenBlock(balance, a.Iban)

	if err := block.Sign(a.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

// CreateReceiveBlock creates a signed ReceiveBlock with the given arguments.
func (a *Account) CreateReceiveBlock(amt primitives.Amount, prev, src primitives.BlockHash) (*primitives.ReceiveBlock, error) {
	block := primitives.NewReceiveBlock(amt, prev, src)

	if err := block.Sign(a.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

// CreateSendBlock creates a signed SendBlock with the given arguments.
func (a *Account) CreateSendBlock(amt primitives.Amount, dest primitives.IBAN, prev primitives.BlockHash) (*primitives.SendBlock, error) {
	block := primitives.NewSendBlock(amt, dest, prev)

	if err := block.Sign(a.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

// Deserialize decodes byte data encoded by gob.
func (a *Account) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(a)
}

// DeserializeJSON decodes JSON data.
func (a *Account) DeseralizeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(a)
}

// Serialize encodes to byte data using gob.
func (a *Account) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(a)
}

// SerializeJSON encodes to JSON data.
func (a *Account) SeralizeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(a)
}
