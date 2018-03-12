package core

import (
	"log"

	"github.com/kookehs/watchmen/primitives"
)

// Request is an instruction for a Node.
type Request struct {
	Account   *Account
	Blueprint *Blueprint
}

// NewRequest returns a pointer to an initialized Request.
func NewRequest(account *Account, blueprint *Blueprint) *Request {
	return &Request{
		Account:   account,
		Blueprint: blueprint,
	}
}

// TODO: Create an unlimited buffered channel or load balance.

// Defines limits and fees for the system
var (
	// Fees
	TransactionFee primitives.Amount = primitives.NewAmount(0.1)

	// Limits
	InputBufferSize int = 100
)

// Node is the structure responsible for carrying out actions on the network.
type Node struct {
	DPoS   *DPoS
	Ledger *Ledger
}

// NewNode returns a pointer to an initialized Node.
func NewNode(dpos *DPoS, ledger *Ledger) *Node {
	return &Node{
		DPoS:   dpos,
		Ledger: ledger,
	}
}

// Process processes the given request taking necessary actions.
func (n *Node) Process(request *Request) (primitives.Block, error) {
	forger, err := n.DPoS.Update(n.Ledger)

	if err != nil {
		return nil, err
	}

	account := request.Account
	blueprint := request.Blueprint
	block, err := n.DPoS.Round.Forge(account, request.Blueprint)

	if err != nil {
		return nil, err
	}

	if err := block.Sign(account.Key.PrivateKey); err != nil {
		return nil, err
	}

	if err := n.Ledger.AppendBlock(block, account.IBAN); err != nil {
		return nil, err
	}

	reward := primitives.NewAmount(0)
	reward.Copy(ForgeReward)

	switch blueprint.Type {
	case primitives.Change:
		reward.Add(reward, VotingFee)
		reward.Add(reward, TransactionFee)
	case primitives.Delegate:
		account.Delegate = true
		account.Share = blueprint.Share
		n.DPoS.Delegates = append(n.DPoS.Delegates, NewDelegate(account))
		reward.Add(reward, DelegateFee)
		reward.Add(reward, TransactionFee)
	case primitives.Open:
	case primitives.Receive:
	case primitives.Send:
		destination := n.Ledger.Accounts[blueprint.Destination.String()]
		prev := n.Ledger.LatestBlock(destination.IBAN)
		blueprint, err := destination.CreateReceiveBlock(blueprint.Amount, &account.Key.PrivateKey.PublicKey, prev, block)

		if err != nil {
			return nil, err
		}

		_, err = n.Process(NewRequest(destination, blueprint))

		if err != nil {
			return nil, err
		}

		reward.Add(reward, TransactionFee)
	default:
		log.Println("Invalid block type")
	}

	n.Payout(forger.Account, reward)
	return block, nil
}

// Payout distributes the block reward to stakeholders evenly according to share.
func (n *Node) Payout(forger *Account, reward primitives.Amount) {
	stakeholders := n.Ledger.Stakeholders(forger.IBAN)

	// Calculate the amount shared.
	share := primitives.NewAmount(forger.Share)
	share.Quo(share, primitives.NewAmount(100))
	share.Mul(reward, share)

	// Calculate the amount to be kept by forger.
	keep := primitives.NewAmount(0)
	keep.Sub(reward, share)

	// Calculate the split.
	ways := primitives.NewAmount(float64(len(stakeholders)))
	reward.Quo(reward, ways)

	for _, stakeholder := range stakeholders {
		if err := n.Transfer(reward, stakeholder, forger); err != nil {
			log.Println(err)
		}
	}

	// TODO: Send forger their split.
}

// Transfer moves funds from one Account to another.
func (n *Node) Transfer(amt primitives.Amount, dst, src *Account) error {
	prev := n.Ledger.LatestBlock(src.IBAN)
	blueprint, err := src.CreateSendBlock(amt, dst.IBAN, prev)

	if err != nil {
		return err
	}

	if _, err := n.Process(NewRequest(src, blueprint)); err != nil {
		return err
	}

	return nil
}
