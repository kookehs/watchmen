package core

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/kookehs/watchmen/primitives"
)

// Defines limits and fees for the system
var (
	// Fees
	DelegateFee primitives.Amount = primitives.NewAmount(25)
	VotingFee   primitives.Amount = primitives.NewAmount(1)

	// Limits
	MaxDelegatesPerAccount int = 101
	MaxDelegatesPerBlock   int = 33
	MaxForgers             int = 101

	// Rewards
	ForgeReward primitives.Amount = primitives.NewAmount(5)
)

// Blueprint contains information used to create a block.
type Blueprint struct {
	Amount      primitives.Amount
	Delegate    bool
	Delegates   []primitives.IBAN
	Destination primitives.IBAN
	Type        primitives.BlockType
	Previous    primitives.Block
	Source      primitives.Block
}

// Delegate contains an Account and their total weight.
type Delegate struct {
	Account *Account
	Weight  primitives.Amount
}

// NewDelegate returns a pointer to a Delegate with a weight of 0.
func NewDelegate(account *Account) *Delegate {
	return &Delegate{
		Account: account,
		Weight:  primitives.NewAmount(0),
	}
}

// TODO: Add sorting functions.

// DPoS contains variables and logic related to the delegated proof of stake.
type DPoS struct {
	// All delegates with their respective total weight
	Delegates []*Delegate
	Round     *Round
}

// NewDPoS returns a pointer to an initialized DPoS.
func NewDPoS() *DPoS {
	return &DPoS{
		Delegates: make([]*Delegate, 0),
		Round:     NewRound(),
	}
}

// CheckMaxDelegateLimit ensures accounts don't vote for more delegates than MaxDelegates.
func CheckMaxDelegateLimit(account *Account, delegates []string) error {
	var add, sub int

	for _, change := range delegates {
		// Change must be atleast 2 characters including symbol
		if len(change) < 2 {
			continue
		}

		switch change[0] {
		case '+':
			add++
		case '-':
			sub++
		default:
			log.Println("Unknown symbol before delegate")
		}
	}

	length := len(account.Delegates) + add - sub

	if length > MaxDelegatesPerAccount {
		return fmt.Errorf("Length of delegates exceeds maximum limit: %v > %v", length, MaxDelegatesPerAccount)
	}

	return nil
}

// ParseDelegateString returns the symbol and Account associated with the given delegate string.
func ParseDelegateString(delegate string, ledger *Ledger) (byte, *Account) {
	symbol := byte(delegate[0])
	username := strings.ToLower(delegate[1:])
	account := ledger.Accounts[username]
	return symbol, account
}

// ParseDelegates updates the delegates for the given Account.
// It may create multiple ChangeBlocks depending on the number of delegates.
func (d *DPoS) ParseDelegates(account *Account, delegates []string, ledger *Ledger, node *Node) ([]*Account, error) {
	length := len(delegates)

	if length == 0 {
		return nil, nil
	}

	split := MaxDelegatesPerBlock

	if length < split {
		split = length
	}

	accounts := make([]*Account, 0)
	elected, err := d.ParseDelegates(account, delegates[split:], ledger, node)

	if err != nil {
		return nil, err
	}

	accounts = append(accounts, elected...)
	ibans := make([]primitives.IBAN, 0)

	for _, change := range delegates[:split] {
		// Change must be atleast 2 characters including symbol
		if len(change) < 2 {
			continue
		}

		symbol, delegate := ParseDelegateString(change, ledger)
		iban := delegate.IBAN
		_, exist := account.Delegates[iban]

		switch symbol {
		case '+':
			if !exist && delegate.Delegate {
				account.Delegates[iban] = true
				accounts = append(accounts, delegate)
				ibans = append(ibans, iban)
			}
		case '-':
			if exist {
				delete(account.Delegates, iban)
				ibans = append(ibans, iban)
			}
		default:
			log.Println("Unknown symbol before delegate")
		}
	}

	if len(ibans) == 0 {
		return nil, nil
	}

	prev := ledger.LatestBlock(account.IBAN)
	blueprint, err := account.CreateChangeBlock(ibans, prev)

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

	return accounts, nil
}

// Elect process the given delegates and distribute the fee to newly elected delegates.
func (d *DPoS) Elect(account *Account, delegates []string, ledger *Ledger, node *Node) error {
	err := CheckMaxDelegateLimit(account, delegates)

	if err != nil {
		return err
	}

	_, err = d.ParseDelegates(account, delegates, ledger, node)

	if err != nil {
		return err
	}

	// TODO: Update delegate weights
	return nil
}

// Round is the system in which each Delegate will forge a single block.
type Round struct {
	Forgers []*Delegate
	Index   int
}

// NewRound returns a pointer to an initialized Round.
func NewRound() *Round {
	return &Round{
		Forgers: make([]*Delegate, MaxForgers),
		Index:   0,
	}
}

// Forge will have the Delegate at Index create the next block.
func (r *Round) Forge(account *Account, blueprint *Blueprint) (primitives.Block, error) {
	var block primitives.Block

	prev := blueprint.Previous
	hash, err := prev.Hash()

	if err != nil {
		return nil, err
	}

	forger := r.Forgers[r.Index]
	r.Index = (r.Index + 1) % MaxForgers

	switch blueprint.Type {
	case primitives.Change:
		block = primitives.NewChangeBlock(prev.Balance(), blueprint.Delegates, hash)
	case primitives.Delegate:
		block = primitives.NewDelegateBlock(prev.Balance(), blueprint.Delegate, hash)
	case primitives.Open:
		block = primitives.NewOpenBlock(blueprint.Amount, account.IBAN)
	case primitives.Receive:
		srcHash, err := blueprint.Source.Hash()

		if err != nil {
			return nil, err
		}

		block = primitives.NewReceiveBlock(blueprint.Amount, hash, srcHash)
	case primitives.Send:
		block = primitives.NewSendBlock(blueprint.Amount, blueprint.Destination, hash)
	default:
		return nil, errors.New("Invalid block type")
	}

	if err := block.SignDelegate(forger.Account.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

// Forger returns the current forger.
func (r *Round) Forger() *Delegate {
	forger := r.Forgers[r.Index]
	return forger
}

func (r *Round) UpdateDelegates(delegates []*Delegate) {
	// TODO: Check if there are new forgers and update
}
