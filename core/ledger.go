package core

import (
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/kookehs/watchmen/primitives"
)

// IBAN is a type alias used for readability and JSON purposes.
type IBAN = string

// Username is a type alias used for readability and JSON purposes.
type Username = string

// Ledger is the structure in which we record accounts and block.
type Ledger struct {
	Accounts map[IBAN]*Account            `json:"accounts"`
	Blocks   map[IBAN][]primitives.Block  `json:"blocks"`
	Users    map[Username]primitives.IBAN `json:"users"`
}

// NewLedger creates and initializes a Ledger for storage of accounts and blocks.
func NewLedger() *Ledger {
	return &Ledger{
		Accounts: make(map[IBAN]*Account),
		Blocks:   make(map[IBAN][]primitives.Block),
		Users:    make(map[Username]primitives.IBAN),
	}
}

// AppendBlock appends the given block to the given IBAN's chain.
func (l *Ledger) AppendBlock(block primitives.Block, iban primitives.IBAN) error {
	if block == nil {
		return errors.New("Cannot append nil block")
	}

	l.Blocks[iban.String()] = append(l.Blocks[iban.String()], block)
	return nil
}

// LatestBlock returns the newest block in the ledger with the given IBAN.
func (l *Ledger) LatestBlock(iban primitives.IBAN) primitives.Block {
	blocks, ok := l.Blocks[iban.String()]

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
	l.Accounts[account.IBAN.String()] = account

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
	l.Accounts[account.IBAN.String()] = account

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

	share := 100.0
	delegate := primitives.NewDelegateBlock(prev.Balance(), hash, share)

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

	account.Delegates[account.IBAN.String()] = true
	return account, nil
}

// OpenGenesisDelegates creates the initial MaxDelegatesPerAccount delegates.
// This receiver should be called after creating the genesis account.
// The process follows the rules of the system returning the delegates.
func (l *Ledger) OpenGenesisDelegates(dpos *DPoS, genesis *Account, node *Node) []*Account {
	delegates := make([]*Account, MaxDelegatesPerAccount)

	split := primitives.NewAmount(0)
	balance := l.LatestBlock(genesis.IBAN).Balance()
	split.Copy(balance)
	fee := primitives.NewAmount(0)
	ways := primitives.NewAmount(float64(MaxDelegatesPerAccount))
	fee.Mul(ways, TransactionFee)
	split.Sub(split, fee)
	split.Quo(split, ways)

	for i := 0; i < MaxDelegatesPerAccount; i++ {
		username := "genesis_" + strconv.Itoa(i+1)
		delegate, err := l.OpenAccount(node, username)

		if err != nil {
			log.Println(err)
			continue
		}

		prev := l.LatestBlock(genesis.IBAN)
		blueprint, err := genesis.CreateSendBlock(split, delegate.IBAN, prev)

		if err != nil {
			log.Println(err)
			continue
		}

		if _, err := node.Process(NewRequest(genesis, blueprint)); err != nil {
			log.Println(err)
			continue
		}

		prev = l.LatestBlock(delegate.IBAN)
		share := 100.0
		blueprint, err = delegate.CreateDelegateBlock(prev, share)

		if err != nil {
			log.Println(err)
			continue
		}

		if _, err = node.Process(NewRequest(delegate, blueprint)); err != nil {
			log.Println(err)
			continue
		}

		if err := dpos.Elect(delegate, []string{"+" + username}, l, node); err != nil {
			log.Println(err)
			continue
		}

		delegates[i] = delegate
	}

	return delegates
}

// Stakeholders returns a list of accounts who elected the given delegate.
func (l *Ledger) Stakeholders(delegate primitives.IBAN) []*Account {
	stakeholders := make([]*Account, 0)

	for _, account := range l.Accounts {
		if _, exist := account.Delegates[delegate.String()]; exist {
			stakeholders = append(stakeholders, account)
		}
	}

	return stakeholders
}

// Username returns the username associated with the given IBAN.
func (l *Ledger) Username(iban primitives.IBAN) string {
	var username string

	for key, value := range l.Users {
		if value == iban {
			username = key
			break
		}
	}

	return username
}

// Deserialize decodes byte data encoded by gob.
func (l *Ledger) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(l)
}

// Deserialize decodes JSON data.
func (l *Ledger) DeserializeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(l)
}

// Serialize encodes to byte data using gob.
func (l *Ledger) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(l)
}

// Serialize encodes to JSON data.
func (l *Ledger) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(l)
}
