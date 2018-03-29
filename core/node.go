package core

import (
	"log"

	"github.com/kookehs/watchmen/primitives"
)

// TODO: Implement network.

// Defines fees for the system
var (
	// Fees
	TransactionFee primitives.Amount = primitives.NewAmount(0.1)
)

// Status interface contains functions related to the statuses of a delegate.
type Status interface {
	Available(string) bool
}

// Node is the structure responsible for carrying out actions on the network.
type Node struct {
	DPoS   *DPoS
	Ledger *Ledger
	Status Status
}

// NewNode returns a pointer to an initialized Node.
func NewNode(dpos *DPoS, ledger *Ledger, status Status) *Node {
	return &Node{
		DPoS:   dpos,
		Ledger: ledger,
		Status: status,
	}
}

// Process processes the given request taking necessary actions.
func (n *Node) Process(request *Request) (primitives.Block, error) {
	var forger *Delegate

	for {
		n.DPoS.Update(n.Ledger)
		forger, err := n.DPoS.Round.Forger()

		if err != nil {
			return nil, err
		}

		if n.Status.Available(n.Ledger.Username(forger.Account.IBAN)) {
			break
		}

		forger.Account.Missed++
		n.DPoS.Round.Index++
	}

	account := request.Account
	blueprint := request.Blueprint
	block, err := n.DPoS.Round.Forge(account, request.Blueprint)

	if err != nil {
		return nil, err
	}

	forger.Account.Forged++

	if err := block.Sign(account.Key.PrivateKey); err != nil {
		return nil, err
	}

	if err := n.Ledger.AppendBlock(block, account.IBAN); err != nil {
		return nil, err
	}

	reward := primitives.NewAmount(0)

	switch blueprint.Type {
	case primitives.Change:
		reward.Copy(ForgeReward)
		reward.Add(reward, VotingFee)
		reward.Add(reward, TransactionFee)
	case primitives.Delegate:
		account.Delegate = true
		account.Share = blueprint.Share
		n.DPoS.Delegates = append(n.DPoS.Delegates, NewDelegate(account))
		reward.Copy(ForgeReward)
		reward.Add(reward, DelegateFee)
		reward.Add(reward, TransactionFee)
	case primitives.Open:
		reward.Copy(ForgeReward)
	case primitives.Receive:
		// No reward for forging a ReceiveBlock.
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

		reward.Copy(ForgeReward)
		reward.Add(reward, TransactionFee)
	default:
		log.Println("Invalid block type")
	}

	if reward.Cmp(primitives.NewAmount(0)) == 1 {
		n.Payout(forger.Account, reward)
	}

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

	if share.Cmp(primitives.NewAmount(0)) == 1 {
		// Calculate the split.
		ways := primitives.NewAmount(float64(len(stakeholders)))
		share.Quo(share, ways)

		for _, stakeholder := range stakeholders {
			prev := n.Ledger.LatestBlock(stakeholder.IBAN)
			blueprint, err := stakeholder.CreateReceiveBlock(share, nil, prev, nil)

			if err != nil {
				log.Println(err)
				continue
			}

			_, err = n.Process(NewRequest(stakeholder, blueprint))

			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

	if keep.Cmp(primitives.NewAmount(0)) == 1 {
		prev := n.Ledger.LatestBlock(forger.IBAN)
		blueprint, err := forger.CreateReceiveBlock(keep, nil, prev, nil)

		if err != nil {
			log.Println(err)
			return
		}

		_, err = n.Process(NewRequest(forger, blueprint))

		if err != nil {
			log.Println(err)
			return
		}
	}
}

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
