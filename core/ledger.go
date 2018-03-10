package core

import (
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/kookehs/watchmen/primitives"
)

// Ledger is the structure in which we record accounts and block.
type Ledger struct {
	Accounts map[primitives.IBAN]*Account           `json:"accounts"`
	Blocks   map[primitives.IBAN][]primitives.Block `json:"blocks"`
	Users    map[string]primitives.IBAN             `json:"users"`
}

// NewLedger creates and initializes a Ledger for storage of accounts and blocks.
func NewLedger() *Ledger {
	return &Ledger{
		Accounts: make(map[primitives.IBAN]*Account),
		Blocks:   make(map[primitives.IBAN][]primitives.Block),
		Users:    make(map[string]primitives.IBAN),
	}
}

// AppendBlock appends the given block to the given IBAN's chain.
func (l *Ledger) AppendBlock(block primitives.Block, iban primitives.IBAN) error {
	if block == nil {
		return errors.New("Cannot append nil block")
	}

	l.Blocks[iban] = append(l.Blocks[iban], block)
	return nil
}

// LatestBlock returns the newest block in the ledger with the given IBAN.
func (l *Ledger) LatestBlock(iban primitives.IBAN) primitives.Block {
	blocks, ok := l.Blocks[iban]

	if ok {
		return blocks[len(blocks)-1]
	}

	return nil
}

// OpenAccount creates an Account for the given username.
func (l *Ledger) OpenAccount(node *Node, username string) (*Account, error) {
	if _, exist := l.Users[username]; exist {
		return nil, fmt.Errorf("Account for %v already exists", username)
	}

	key, err := primitives.NewKeyForICAP(rand.Reader)

	if err != nil {
		return nil, err
	}

	account := NewAccount(key)
	amount := primitives.NewAmount(0)
	blueprint, err := account.CreateOpenBlock(amount)

	if err != nil {
		return nil, err
	}

	username = strings.ToLower(username)
	l.Users[username] = account.IBAN
	l.Accounts[account.IBAN] = account

	request := NewRequest(account, blueprint)

	if _, err := node.Process(request); err != nil {
		return nil, err
	}

	return account, nil
}

// OpenGenesisAccount creates an initial account that bypasses the system.
// Creates an account with an initial amount with delegate status.
// This method is meant to be called once to initialize the system.
func (l *Ledger) OpenGenesisAccount(username string) (*Account, error) {
	if _, exist := l.Users[username]; exist {
		return nil, fmt.Errorf("Account for %v already exists", username)
	}

	key, err := primitives.NewKeyForICAP(rand.Reader)

	if err != nil {
		return nil, err
	}

	account := NewAccount(key)
	amount := primitives.NewAmount(100000000)
	open := primitives.NewOpenBlock(amount, account.IBAN)

	if open == nil {
		return nil, errors.New("Unable to create block")
	}

	username = strings.ToLower(username)
	l.Users[username] = account.IBAN
	l.Accounts[account.IBAN] = account

	if err := open.Sign(account.Key.PrivateKey); err != nil {
		return nil, err
	}

	// Self-sign the opening block for the genesis account.
	if err := open.SignWitness(account.Key.PrivateKey); err != nil {
		return nil, err
	}

	if err := l.AppendBlock(open, account.IBAN); err != nil {
		return nil, err
	}

	prev := l.LatestBlock(account.IBAN)
	hash, err := prev.Hash()

	if err != nil {
		return nil, err
	}

	delegate := primitives.NewDelegateBlock(prev.Balance(), true, hash)

	if delegate == nil {
		return nil, errors.New("Unable to create block")
	}

	if err := delegate.Sign(account.Key.PrivateKey); err != nil {
		return nil, err
	}

	// Self-sign the opening block for the genesis account.
	if err := delegate.SignWitness(account.Key.PrivateKey); err != nil {
		return nil, err
	}

	if err := l.AppendBlock(delegate, account.IBAN); err != nil {
		return nil, err
	}

	account.Delegate = true
	prev = l.LatestBlock(account.IBAN)
	hash, err = prev.Hash()

	if err != nil {
		return nil, err
	}

	change := primitives.NewChangeBlock(prev.Balance(), []primitives.IBAN{account.IBAN}, hash)

	if change == nil {
		return nil, errors.New("Unable to create block")
	}

	if err := change.Sign(account.Key.PrivateKey); err != nil {
		return nil, err
	}

	// Self-sign the opening block for the genesis account.
	if err := change.SignWitness(account.Key.PrivateKey); err != nil {
		return nil, err
	}

	if err := l.AppendBlock(change, account.IBAN); err != nil {
		return nil, err
	}

	account.Delegates[account.IBAN] = true
	return account, nil
}

// Deserialize decodes byte data encoded by gob.
func (l *Ledger) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(l)
}

// Deserialize decodes json data.
func (l *Ledger) DeserializeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(l)
}

// Serialize encodes to byte data using gob.
func (l *Ledger) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(l)
}

// Serialize encodes to json data.
func (l *Ledger) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(l)
}
