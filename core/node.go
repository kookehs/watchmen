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

	switch blueprint.Type {
	case primitives.Change:
		// Transfer of voting fee
		if err := n.Transfer(VotingFee, forger.Account, account); err != nil {
			log.Println(err)
		}
	case primitives.Delegate:
		// Transfer of delegate fee
		if err := n.Transfer(DelegateFee, forger.Account, account); err != nil {
			log.Println(err)
			break
		}

		account.Delegate = true
		n.DPoS.Delegates = append(n.DPoS.Delegates, NewDelegate(account))
	case primitives.Open:
	case primitives.Receive:
	case primitives.Send:
		destination := n.Ledger.Accounts[blueprint.Destination]
		prev := n.Ledger.LatestBlock(destination.IBAN)
		blueprint, err := destination.CreateReceiveBlock(blueprint.Amount, &account.Key.PrivateKey.PublicKey, prev, block)

		if err != nil {
			return nil, err
		}

		_, err = n.Process(NewRequest(destination, blueprint))

		if err != nil {
			return nil, err
		}
	default:
		log.Println("Invalid block type")
	}

	// TODO: Payout fees + reward
	return block, nil
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
