package core

import (
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/kookehs/watchmen/primitives"
)

// Ledger is the structure in which we record accounts and block.
type Ledger struct {
	// Mapping of usernames to Account
	Accounts map[string]*Account                    `json:"accounts"`
	Blocks   map[primitives.IBAN][]primitives.Block `json:"blocks"`
}

// NewLedger creates and initializes a Ledger for storage of accounts and blocks.
func NewLedger() *Ledger {
	return &Ledger{
		Accounts: make(map[string]*Account),
		Blocks:   make(map[primitives.IBAN][]primitives.Block),
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
	key, err := primitives.NewKeyForICAP(rand.Reader)

	if err != nil {
		return nil, err
	}

	account := NewAccount(key)
	username = strings.ToLower(username)
	l.Accounts[username] = account
	blueprint, err := account.CreateOpenBlock()

	if err != nil {
		return nil, err
	}

	output := make(chan primitives.Block)
	node.Input <- Request{Account: account, Blueprint: blueprint, Output: output}
	block := <-output
	close(output)

	if block == nil {
		return nil, errors.New("Unable to forge block")
	}

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
