package core

import (
	"crypto/rand"

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
