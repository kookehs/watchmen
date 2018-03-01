package core

import (
	"fmt"
	"log"

	"github.com/kookehs/watchmen/primitives"
)

// Defines limits and fees for the system
var (
	MaxDelegates         = 101
	MaxDelegatesPerBlock = 33
	VotingFee            = 1
)

// DPoS contains variables and logic related to the delegated proof of stake.
type DPoS struct {
	// All delegates with respective total weight
	Delegates map[primitives.IBAN]uint64
}

// Delegate updates the delegates for the given Account.
// It may create multiple ChangeBlocks depending on the number of delegates.
func Delegate(account *Account, delegates []string, ledger *Ledger) error {
	length := len(delegates)

	if length == 0 {
		return nil
	}

	if length > MaxDelegates {
		return fmt.Errorf("Length of delegates exceeds maximum limit: %v > %v", length, MaxDelegates)
	}

	split := MaxDelegatesPerBlock

	if length < split {
		split = length
	}

	Delegate(account, delegates[split:], ledger)
	changes := make([]primitives.IBAN, 0)

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
				// TODO: Distribute voting fee to new delegates.
				account.Delegates[iban] = true
				changes = append(changes, iban)
			}
		case '-':
			if exist {
				delete(account.Delegates, iban)
				changes = append(changes, iban)
			}
		default:
			log.Println("Unknown symbol before delegate")
		}
	}

	if len(changes) == 0 {
		return nil
	}

	prev := ledger.LatestBlock(account.IBAN)
	hash, err := prev.Hash()

	if err != nil {
		return err
	}

	block, err := account.CreateChangeBlock(prev.Balance(), changes, hash)

	if err != nil {
		return err
	}

	if err := ledger.AppendBlock(block, account.IBAN); err != nil {
		return err
	}

	return nil
}

// ParseDelegateString returns the symbol and IBAN associated with the given delegate string.
func ParseDelegateString(delegate string, ledger *Ledger) (byte, *Account) {
	symbol := byte(delegate[0])
	username := delegate[1:]
	account := ledger.Accounts[username]
	return symbol, account
}
