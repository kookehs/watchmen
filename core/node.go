package core

import (
	"log"

	"github.com/kookehs/watchmen/primitives"
)

// Request is an instruction for a Node.
type Request struct {
	Account   *Account
	Blueprint *Blueprint
	Output    chan primitives.Block
}

// NewRequest returns a pointer to an initialized Request.
func NewRequest(account *Account, blueprint *Blueprint, output chan primitives.Block) *Request {
	return &Request{
		Account:   account,
		Blueprint: blueprint,
		Output:    output,
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
	Input chan *Request
}

// NewNode returns a pointer to an initialized Node.
func NewNode() *Node {
	return &Node{
		Input: make(chan *Request, InputBufferSize),
	}
}

// Listen waits for requests coming through Input.
func (n *Node) Listen(dpos *DPoS, ledger *Ledger) {
	// TODO: Who forges when there are no delegates?
	// TODO: Generate genesis delegates.
	for {
		request := <-n.Input
		forger := dpos.Update(ledger).Account
		account := request.Account
		failed := false
		block, err := dpos.Round.Forge(account, request.Blueprint)

		if err != nil {
			failed = true
			log.Println(err)
		} else if err := block.Sign(account.Key.PrivateKey); err != nil {
			failed = true
			log.Println(err)
		} else if err := ledger.AppendBlock(block, account.IBAN); err != nil {
			failed = true
			log.Println(err)
		}

		if !failed {
			switch request.Blueprint.Type {
			case primitives.Change:
				// Transfer of voting fee
				if err := n.Transfer(VotingFee, forger, account, ledger); err != nil {
					log.Println(err)
				}
			case primitives.Delegate:
				// Transfer of delegate fee
				if err := n.Transfer(DelegateFee, forger, account, ledger); err != nil {
					log.Println(err)
					break
				}

				account.Delegate = true
				dpos.Delegates = append(dpos.Delegates, NewDelegate(account))
			case primitives.Open:
			case primitives.Receive:
			case primitives.Send:
			default:
				log.Println("Invalid block type")
			}

			// TODO: Payout fees + reward
		}

		if request.Output != nil {
			request.Output <- block
		}
	}
}

// Transfer moves funds from one Account to another creating
// the necessary send and receive blocks.
func (n *Node) Transfer(amt primitives.Amount, dst, src *Account, ledger *Ledger) error {
	blueprint, err := src.CreateSendBlock(amt, src.IBAN, ledger.LatestBlock(src.IBAN))

	if err != nil {
		return err
	}

	output := make(chan primitives.Block)
	n.Input <- NewRequest(src, blueprint, output)
	block := <-output
	close(output)
	blueprint, err = dst.CreateReceiveBlock(amt, &src.Key.PrivateKey.PublicKey, ledger.LatestBlock(dst.IBAN), block)

	if err != nil {
		return err
	}

	n.Input <- NewRequest(dst, blueprint, output)
	return nil
}
