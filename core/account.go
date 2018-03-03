package core

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"

	"github.com/kookehs/watchmen/primitives"
)

// Account contains address as well as the key that generated it.
type Account struct {
	BBAN      primitives.BBAN          `json:"bban"`
	Delegate  bool                     `json:"delegate"`
	Delegates map[primitives.IBAN]bool `json:"delegates"`
	IBAN      primitives.IBAN          `json:"iban"`
	Key       *primitives.Key          `json:"key"`
}

// NewAccount creates and initializes an account with the given key.
func NewAccount(key *primitives.Key) *Account {
	bban := primitives.MakeBBAN([]byte(key.Address.String()))
	iban := primitives.MakeIBAN([]byte("TV00" + bban.String()))

	return &Account{
		BBAN:      bban,
		Delegate:  false,
		Delegates: make(map[primitives.IBAN]bool),
		IBAN:      iban,
		Key:       key,
	}
}

// CreateChangeBlock creates a signed ChangeBlock with the given arguments.
func (a *Account) CreateChangeBlock(delegates []primitives.IBAN, prev primitives.Block) (*primitives.ChangeBlock, error) {
	hash, err := prev.Hash()

	if err != nil {
		return nil, err
	}

	block := primitives.NewChangeBlock(prev.Balance(), delegates, hash)

	if err := block.Sign(a.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

// CreateDelegateBlock creates a signed DelegateBlock with the given arguments.
func (a *Account) CreateDelegateBlock(delegate bool, prev primitives.Block) (*primitives.DelegateBlock, error) {
	hash, err := prev.Hash()

	if err != nil {
		return nil, err
	}

	block := primitives.NewDelegateBlock(prev.Balance(), delegate, hash)

	if err := block.Sign(a.Key.PrivateKey); err != nil {
		return nil, err
	}

	a.Delegate = true
	return block, nil
}

// CreateOpenBlock creates a signed OpenBlock.
func (a *Account) CreateOpenBlock() (*primitives.OpenBlock, error) {
	balance := primitives.NewAmount(100)
	block := primitives.NewOpenBlock(balance, a.IBAN)

	if err := block.Sign(a.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

// CreateReceiveBlock creates a signed ReceiveBlock with the given arguments.
func (a *Account) CreateReceiveBlock(amt primitives.Amount, prev primitives.Block, src primitives.BlockHash) (*primitives.ReceiveBlock, error) {
	hash, err := prev.Hash()

	if err != nil {
		return nil, err
	}

	balance := primitives.NewAmount(0)
	balance.Add(prev.Balance(), amt)
	block := primitives.NewReceiveBlock(balance, hash, src)

	if err := block.Sign(a.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

// CreateSendBlock creates a signed SendBlock with the given arguments.
func (a *Account) CreateSendBlock(amt primitives.Amount, dest primitives.IBAN, prev primitives.Block) (*primitives.SendBlock, error) {
	if prev.Balance().Cmp(amt) == -1 {
		return nil, errors.New("Insufficient funds")
	}

	hash, err := prev.Hash()

	if err != nil {
		return nil, err
	}

	balance := primitives.NewAmount(0)
	balance.Sub(prev.Balance(), amt)
	block := primitives.NewSendBlock(balance, dest, hash)

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
