package core

import (
	"fmt"
	"log"
	"strings"

	"github.com/kookehs/watchmen/primitives"
)

// Defines limits and fees for the system
var (
	// Limits
	MaxDelegatesPerAccount int = 101
	MaxDelegatesPerBlock   int = 33
	MaxForgers             int = 101

	// Fees
	DelegateFee          primitives.Amount = primitives.NewAmount(25)
	MinimumDelegateStake primitives.Amount = primitives.NewAmount(25)
	VotingFee            primitives.Amount = primitives.NewAmount(1)
)

// Delegate contains an IBAN and their total weight.
type Delegate struct {
	IBAN   primitives.IBAN
	Weight primitives.Amount
}

// NewDelegate returns a pointer to a Delegate with a weight of 0.
func NewDelegate(iban primitives.IBAN) *Delegate {
	return &Delegate{
		IBAN:   iban,
		Weight: primitives.NewAmount(0),
	}
}

// TODO: Add sorting functions.

// DPoS contains variables and logic related to the delegated proof of stake.
type DPoS struct {
	// All delegates with their respective total weight
	Delegates []*Delegate
	Forger    uint64
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

// ParseDelegates updates the delegates for the given Account.
// It may create multiple ChangeBlocks depending on the number of delegates.
func ParseDelegates(account *Account, delegates []string, ledger *Ledger) ([]*Account, error) {
	length := len(delegates)

	if length == 0 {
		return nil, nil
	}

	split := MaxDelegatesPerBlock

	if length < split {
		split = length
	}

	accounts := make([]*Account, 0)
	elected, err := ParseDelegates(account, delegates[split:], ledger)

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
	block, err := account.CreateChangeBlock(ibans, prev)

	if err != nil {
		return nil, err
	}

	if err := ledger.AppendBlock(block, account.IBAN); err != nil {
		return nil, err
	}

	return accounts, nil
}

// ParseDelegateString returns the symbol and Account associated with the given delegate string.
func ParseDelegateString(delegate string, ledger *Ledger) (byte, *Account) {
	symbol := byte(delegate[0])
	username := strings.ToLower(delegate[1:])
	account := ledger.Accounts[username]
	return symbol, account
}

// Delegate process the given delegates and distribute the fee to newly elected delegates.
func (d *DPoS) Delegate(account *Account, delegates []string, ledger *Ledger) error {
	err := CheckMaxDelegateLimit(account, delegates)

	if err != nil {
		return err
	}

	_, err = ParseDelegates(account, delegates, ledger)

	if err != nil {
		return err
	}

	// TODO: Voting fee should be given to the current forger
	// TODO: Update delegate weights
	return nil
}

// TODO: Generate a reward after each round and distribute to
// forger and their supporters.

// Round is the system in which each Delegate will forge a single block.
type Round struct {
	Delegate []*Delegate
	Index    uint64
}

func (r *Round) Forge() {

}
