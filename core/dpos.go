package core

import (
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
	// TODO: Split longer delegates into several ChangeBlocks.
	changes := make([]primitives.IBAN, MaxDelegatesPerBlock)

	for _, change := range delegates {
		// Change must be atleast 2 characters including symbol
		if len(change) < 2 {
			continue
		}

		log.Println(delegates)

		symbol := change[0]
		username := change[1:]
		iban := ledger.Accounts[username].IBAN
		_, exist := account.Delegates[iban]

		switch symbol {
		case '+':
			if !exist {
				// TODO: Distribute voting fee to new delegates.
				account.Delegates[iban] = true
			}

			changes = append(changes, iban)
		case '-':
			if exist {
				delete(account.Delegates, iban)
			}

			changes = append(changes, iban)
		default:
			log.Println("Unknown symbol before delegate")
		}
	}

	prev := ledger.LatestBlock(account.IBAN)
	hash, err := prev.Hash()

	if err != nil {
		return err
	}

	// TODO: A root function of CreateChangeBlock should call this function instead.
	account.CreateChangeBlock(prev.Balance(), changes, hash)
	return nil
}
