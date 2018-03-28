package core

import (
	"errors"
	"fmt"
	"log"
	"sort"
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
	// TODO: Consider decreasing reward amount.
	ForgeReward primitives.Amount = primitives.NewAmount(4)
)

// Blueprint contains information used to create a block.
type Blueprint struct {
	Amount      primitives.Amount
	Balance     primitives.Amount
	Delegates   []primitives.IBAN
	Destination primitives.IBAN
	Previous    primitives.Block
	Share       float64
	Source      primitives.Block
	Type        primitives.BlockType
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

// Delegates is a wrapper type for sorting
type Delegates []*Delegate

// Len returns the length of Delegates.
func (d Delegates) Len() int {
	return len(d)
}

// Swap swaps delegates at position i and j.
func (d Delegates) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less returns if the weight at i is less than j.
// Sort in descending order because we want those with highest weights.
func (d Delegates) Less(i, j int) bool {
	return (d[i].Weight.Cmp(d[j].Weight) == 1)
}

// DPoS contains variables and logic related to the delegated proof of stake.
type DPoS struct {
	// All delegates with their respective total weight
	Delegates Delegates
	Round     *Round
}

// NewDPoS returns a pointer to an initialized DPoS.
func NewDPoS() *DPoS {
	return &DPoS{
		Delegates: make(Delegates, 0),
		Round:     NewRound(Delegates{}),
	}
}

// CalculateWeights iterates through all accounts and their delegates.
// The final weight is the sum of the amounts each supporter holds.
func CalculateWeights(ledger *Ledger) Delegates {
	// Map for quick look up. Slice for sorting.
	delegates := make(map[IBAN]*Delegate)
	values := make(Delegates, 0)

	for _, account := range ledger.Accounts {
		prev := ledger.LatestBlock(account.IBAN)

		if prev == nil {
			continue
		}

		weight := prev.Balance()

		for iban, _ := range account.Delegates {
			if _, exist := delegates[iban]; !exist {
				elected := NewDelegate(ledger.Accounts[iban])
				delegates[iban] = elected
				values = append(values, elected)
			}

			delegates[iban].Weight.Add(delegates[iban].Weight, weight)
		}
	}

	sort.Sort(values)
	return values
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

// Elect processes the given delegates and distribute the fee to newly elected delegates.
func (d *DPoS) Elect(account *Account, delegates []string, ledger *Ledger, node *Node) error {
	err := CheckMaxDelegateLimit(account, delegates)

	if err != nil {
		return err
	}

	_, err = d.ParseDelegates(account, delegates, ledger, node)

	if err != nil {
		return err
	}

	return nil
}

// ParseDelegateString returns the symbol and Account associated with the given delegate string.
func ParseDelegateString(delegate string, ledger *Ledger) (byte, *Account) {
	symbol := byte(delegate[0])
	username := strings.ToLower(delegate[1:])
	iban := ledger.Users[username]
	account := ledger.Accounts[iban.String()]
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
		_, exist := account.Delegates[iban.String()]

		switch symbol {
		case '+':
			if !exist && delegate.Delegate {
				account.Delegates[iban.String()] = true
				accounts = append(accounts, delegate)
				ibans = append(ibans, iban)
			}
		case '-':
			if exist {
				delete(account.Delegates, iban.String())
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

	request := NewRequest(account, blueprint)

	if _, err := node.Process(request); err != nil {
		return nil, err
	}

	return accounts, nil
}

// Update checks if a new round needs to be created and returns the current forger.
func (d *DPoS) Update(ledger *Ledger) Delegates {
	if (len(d.Round.Forgers) == 0) || ((d.Round.Index != 0) && ((d.Round.Index % len(d.Round.Forgers)) == 0)) {
		d.Delegates = CalculateWeights(ledger)
		d.Round = NewRound(d.Delegates)
	}

	return d.Round.Forgers
}

// Round is the system in which each Delegate will forge a single block.
type Round struct {
	Forgers Delegates
	Index   int
}

// NewRound returns a pointer to an initialized Round.
func NewRound(delegates Delegates) *Round {
	split := MaxForgers

	if len(delegates) < split {
		split = len(delegates)
	}

	return &Round{
		Forgers: delegates[:split],
		Index:   0,
	}
}

// Forge will have the Delegate at Index create the next block.
func (r *Round) Forge(account *Account, blueprint *Blueprint) (primitives.Block, error) {
	var block, prev primitives.Block
	var hash primitives.BlockHash
	var err error

	if blueprint.Type != primitives.Open {
		prev = blueprint.Previous

		hash, err = prev.Hash()

		if err != nil {
			return nil, err
		}
	}

	forger := r.Forgers[r.Index]
	r.Index++

	switch blueprint.Type {
	case primitives.Change:
		block = primitives.NewChangeBlock(blueprint.Balance, blueprint.Delegates, hash)
	case primitives.Delegate:
		block = primitives.NewDelegateBlock(blueprint.Balance, hash, blueprint.Share)
	case primitives.Open:
		block = primitives.NewOpenBlock(blueprint.Balance, account.IBAN)
	case primitives.Receive:
		srcHash := primitives.BlockHashZero

		if blueprint.Source != nil {
			if srcHash, err = blueprint.Source.Hash(); err != nil {
				return nil, err
			}
		}

		block = primitives.NewReceiveBlock(blueprint.Balance, hash, srcHash)
	case primitives.Send:
		block = primitives.NewSendBlock(blueprint.Balance, blueprint.Destination, hash)
	default:
		return nil, errors.New("Invalid block type")
	}

	if block == nil {
		return nil, errors.New("Unabled to forge block")
	}

	if err := block.SignWitness(forger.Account.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

// Forger returns the current forger.
func (r *Round) Forger() (*Delegate, error) {
	if (len(r.Forgers) == 0) || ((r.Index != 0) && ((r.Index % len(r.Forgers)) == 0)) {
		return nil, errors.New("No current forger. Round has ended.")
	}

	forger := r.Forgers[r.Index]
	return forger, nil
}
