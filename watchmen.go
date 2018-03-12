package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"log"

	"github.com/kookehs/watchmen/core"
	"github.com/kookehs/watchmen/primitives"
)

func init() {
	gob.Register(elliptic.P256())
	amount := primitives.NewAmount(0)
	hash := primitives.BlockHashZero
	iban := primitives.IBAN{}
	gob.Register(primitives.NewChangeBlock(amount, []primitives.IBAN{}, hash))
	gob.Register(primitives.NewDelegateBlock(amount, hash, 0))
	gob.Register(primitives.NewOpenBlock(amount, iban))
	gob.Register(primitives.NewReceiveBlock(amount, hash, hash))
	gob.Register(primitives.NewSendBlock(amount, iban, hash))
}

func main() {
	// TODO: Implement network for blockchain.

	// Move the below to respective test files.
	ledger := core.NewLedger()

	genesis, err := ledger.OpenGenesisAccount("Genesis Account")

	if err != nil {
		log.Println(err)
		return
	}

	dpos := core.NewDPoS()
	node := core.NewNode(dpos, ledger)
	ledger.OpenGenesisDelegates(dpos, genesis, node)

	var buffer bytes.Buffer

	if err := ledger.SerializeJSON(&buffer); err != nil {
		log.Println(err)
		return
	}

	log.Println(buffer.String())
}
