package core

import (
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	"io"

	"github.com/kookehs/watchmen/primitives"
)

// Ledger is the structure in which we record accounts and block.
type Ledger struct {
	// Mapping of usernames to Account
	Accounts map[string]Account
	// Mapping of IBAN to []Block
	Blocks map[primitives.IBAN][]primitives.Block
}

// NewLedger creates and initializes a Ledger for storage of accounts and blocks.
func NewLedger() *Ledger {
	return &Ledger{
		Accounts: make(map[string]Account),
		Blocks:   make(map[primitives.IBAN][]primitives.Block),
	}
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
func (l *Ledger) OpenAccount(username string) error {
	key, err := primitives.NewKeyForICAP(rand.Reader)

	if err != nil {
		return err
	}

	account := MakeAccount(*key)
	l.Accounts[username] = account
	block, err := account.CreateOpenBlock()

	if err != nil {
		return err
	}

	chain := l.Blocks[account.IBAN]
	chain = append(chain, block)
	return nil
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
