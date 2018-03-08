package core

import (
	"crypto/ecdsa"
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

// VerifyBlock returns whether or not the block was signed by the given key.
func VerifyBlock(block primitives.Block, key *ecdsa.PublicKey) error {
	verified, err := block.Verify(key)

	if err != nil {
		return err
	}

	if !verified {
		return errors.New("Block was not signed with the given key")
	}

	return nil
}

// CreateChangeBlock creates a blueprint for a ChangeBlock with the given arguments.
func (a *Account) CreateChangeBlock(delegates []primitives.IBAN, prev primitives.Block) (*Blueprint, error) {
	if err := VerifyBlock(prev, &a.Key.PrivateKey.PublicKey); err != nil {
		return nil, err
	}

	blueprint := &Blueprint{
		Delegates: delegates,
		Previous:  prev,
		Type:      primitives.Change,
	}

	return blueprint, nil
}

// CreateDelegateBlock creates a blueprint for a DelegateBlock with the given arguments.
func (a *Account) CreateDelegateBlock(delegate bool, prev primitives.Block) (*Blueprint, error) {
	if err := VerifyBlock(prev, &a.Key.PrivateKey.PublicKey); err != nil {
		return nil, err
	}

	if prev.Balance().Cmp(DelegateFee) == -1 {
		return nil, errors.New("Insufficient funds")
	}

	if a.Delegate {
		return nil, errors.New("Account is already a delegate")
	}

	// TODO: Allow a % share to supporters.
	blueprint := &Blueprint{
		Delegate: true,
		Previous: prev,
		Type:     primitives.Delegate,
	}

	return blueprint, nil
}

// CreateOpenBlock creates a blueprint for an OpenBlock.
func (a *Account) CreateOpenBlock(amt primitives.Amount) (*Blueprint, error) {
	blueprint := &Blueprint{
		Amount: amt,
		Type:   primitives.Open,
	}

	return blueprint, nil
}

// CreateReceiveBlock creates a blueprint for a ReceiveBlock with the given arguments.
func (a *Account) CreateReceiveBlock(amt primitives.Amount, key *ecdsa.PublicKey, prev, src primitives.Block) (*Blueprint, error) {
	if err := VerifyBlock(prev, &a.Key.PrivateKey.PublicKey); err != nil {
		return nil, err
	}

	if err := VerifyBlock(src, key); err != nil {
		return nil, err
	}

	balance := primitives.NewAmount(0)
	balance.Add(prev.Balance(), amt)

	blueprint := &Blueprint{
		Amount:   balance,
		Previous: prev,
		Source:   src,
		Type:     primitives.Receive,
	}

	return blueprint, nil
}

// CreateSendBlock creates a blueprint for a SendBlock with the given arguments.
func (a *Account) CreateSendBlock(amt primitives.Amount, dst primitives.IBAN, prev primitives.Block) (*Blueprint, error) {
	if err := VerifyBlock(prev, &a.Key.PrivateKey.PublicKey); err != nil {
		return nil, err
	}

	if prev.Balance().Cmp(amt) == -1 {
		return nil, errors.New("Insufficient funds")
	}

	balance := primitives.NewAmount(0)
	balance.Sub(prev.Balance(), amt)

	blueprint := &Blueprint{
		Amount:      balance,
		Destination: dst,
		Previous:    prev,
		Type:        primitives.Send,
	}

	return blueprint, nil
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
