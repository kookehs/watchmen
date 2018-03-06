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

// TODO: Create an unlimited buffered channel or load balance.

// Limits for the network.
var (
	InputBufferSize = 100
)

// Node is the structure responsible for carrying out actions on the network.
type Node struct {
	Input chan Request
}

// NewNode returns a pointer to an initialized Node.
func NewNode() *Node {
	return &Node{
		Input: make(chan Request, InputBufferSize),
	}
}

// Listen waits for requests coming through Input.
func (n *Node) Listen(dpos *DPoS, ledger *Ledger) {
	for {
		request := <-n.Input
		forger := dpos.Round.Forger().Account
		account := request.Account
		block, err := dpos.Round.Forge(account, request.Blueprint)

		if err != nil {
			log.Println(err)
			continue
		}

		if err := block.Sign(account.Key.PrivateKey); err != nil {
			log.Println(err)
			continue
		}

		if err := ledger.AppendBlock(block, account.IBAN); err != nil {
			log.Println(err)
			continue
		}

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
			}
		case primitives.Open:
		case primitives.Receive:
		case primitives.Send:
		default:
			log.Println("Invalid block type")
		}

		// TODO: Payout reward

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
	n.Input <- Request{Account: src, Blueprint: blueprint, Output: output}
	block := <-output
	close(output)
	blueprint, err = dst.CreateReceiveBlock(amt, &src.Key.PrivateKey.PublicKey, ledger.LatestBlock(dst.IBAN), block)

	if err != nil {
		return err
	}

	n.Input <- Request{Account: dst, Blueprint: blueprint, Output: nil}
	return nil
}
