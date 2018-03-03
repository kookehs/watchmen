package core

import (
	"fmt"
	"log"

	"github.com/kookehs/watchmen/primitives"
)

// Defines limits and fees for the system
var (
	MaxDelegates         int     = 101
	MaxDelegatesPerBlock int     = 33
	VotingFee            float64 = 1
)

// DPoS contains variables and logic related to the delegated proof of stake.
type DPoS struct {
	// All delegates with respective total weight
	Delegates map[primitives.IBAN]uint64
}

// Delegate process the given delegates and distribute the fee to newly elected delegates.
func Delegate(account *Account, delegates []string, ledger *Ledger) error {
	elected, err := ProcessDelegates(account, delegates, ledger)

	if err != nil {
		return err
	}

	// TODO: Who gets voting fee when only delegates are removed?
	// TODO: Voting fee should be given to the current forger?
	if err := DistributeFee(account, elected, VotingFee, ledger); err != nil {
		return err
	}

	return nil
}

// DistributeFee distributes the voting fee among the newly elected delegates.
func DistributeFee(account *Account, delegates []*Account, funds float64, ledger *Ledger) error {
	split := primitives.NewAmount(funds / float64(len(delegates)))

	for _, delegate := range delegates {
		prev := ledger.LatestBlock(account.IBAN)
		sendBlock, err := account.CreateSendBlock(split, delegate.IBAN, prev)

		if err != nil {
			return err
		}

		if err := ledger.AppendBlock(sendBlock, account.IBAN); err != nil {
			return err
		}

		hash, err := sendBlock.Hash()

		if err != nil {
			return err
		}

		prev = ledger.LatestBlock(delegate.IBAN)
		receiveBlock, err := delegate.CreateReceiveBlock(split, prev, hash)

		if err != nil {
			return err
		}

		if err := ledger.AppendBlock(receiveBlock, delegate.IBAN); err != nil {
			return err
		}
	}

	return nil
}

// ParseDelegateString returns the symbol and Account associated with the given delegate string.
func ParseDelegateString(delegate string, ledger *Ledger) (byte, *Account) {
	symbol := byte(delegate[0])
	username := delegate[1:]
	account := ledger.Accounts[username]
	return symbol, account
}

// ProcessDelegates updates the delegates for the given Account.
// It may create multiple ChangeBlocks depending on the number of delegates.
func ProcessDelegates(account *Account, delegates []string, ledger *Ledger) ([]*Account, error) {
	length := len(delegates)

	if length == 0 {
		return nil, nil
	}

	if length > MaxDelegates {
		return nil, fmt.Errorf("Length of delegates exceeds maximum limit: %v > %v", length, MaxDelegates)
	}

	split := MaxDelegatesPerBlock

	if length < split {
		split = length
	}

	accounts := make([]*Account, 0)
	elected, err := ProcessDelegates(account, delegates[split:], ledger)

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
