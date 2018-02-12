package core

import (
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	"io"

	"github.com/kookehs/watchmen/primitives"
)

type Ledger struct {
	// Mapping of usernames to Account
	Accounts map[string]Account
	// Mapping of IBAN to []Block
	Blocks map[string][]primitives.Block
}

func NewLedger() *Ledger {
	return &Ledger{
		Accounts: make(map[string]Account),
		Blocks:   make(map[string][]primitives.Block),
	}
}

func (l *Ledger) CreateAccount(s string) error {
	key, err := primitives.NewKeyForICAP(rand.Reader)

	if err != nil {
		return err
	}

	account := MakeAccount(*key)
	l.Accounts[s] = account
	block, err := account.CreateOpenBlock()

	if err != nil {
		return err
	}

	chain := l.Blocks[account.Iban.String()]
	chain = append(chain, block)
	return nil
}

func (l *Ledger) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(l); err != nil {
		return err
	}

	return nil
}

func (l *Ledger) DeserializeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(l); err != nil {
		return err
	}

	return nil
}

func (l *Ledger) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(l); err != nil {
		return err
	}

	return nil
}

func (l *Ledger) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(l); err != nil {
		return err
	}

	return nil
}
