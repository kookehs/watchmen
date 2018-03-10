package main

import (
	"crypto/elliptic"
	"encoding/gob"
	"log"
	"strconv"

	"github.com/kookehs/watchmen/core"
	"github.com/kookehs/watchmen/primitives"
)

func main() {
	gob.Register(elliptic.P256())
	// TODO: Implement network for blockchain.

	// Move the below to respective test files.
	ledger := core.NewLedger()

	genesis, err := ledger.OpenGenesisAccount("Genesis Account")

	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("IBAN: %v", genesis.IBAN.String())
	log.Printf("Amount: %v", ledger.LatestBlock(genesis.IBAN).Balance())

	dpos := core.NewDPoS()
	node := core.NewNode(dpos, ledger)

	delegates := make([]string, 101)

	split := primitives.NewAmount(0)
	balance := ledger.LatestBlock(genesis.IBAN).Balance()
	split = split.Copy(balance)
	split.Sub(split, core.TransactionFee)
	ways := primitives.NewAmount(float64(core.MaxDelegatesPerAccount))
	split.Quo(split, ways)

	// TODO: Create a function to generate genesis delegates.
	for i := 0; i < core.MaxDelegatesPerAccount; i++ {
		username := "genesis_" + strconv.Itoa(i+1)
		delegate, err := ledger.OpenAccount(node, username)

		if err != nil {
			log.Println(err)
			continue
		}

		prev := ledger.LatestBlock(genesis.IBAN)
		blueprint, err := genesis.CreateSendBlock(split, delegate.IBAN, prev)

		if err != nil {
			log.Println(err)
			continue
		}

		if _, err := node.Process(core.NewRequest(genesis, blueprint)); err != nil {
			log.Println(err)
			continue
		}

		prev = ledger.LatestBlock(delegate.IBAN)
		blueprint, err = delegate.CreateDelegateBlock(true, prev)

		if err != nil {
			log.Println(err)
			continue
		}

		if _, err = node.Process(core.NewRequest(delegate, blueprint)); err != nil {
			log.Println(err)
			continue
		}

		if err := dpos.Elect(delegate, []string{"+" + username}, ledger, node); err != nil {
			log.Println(err)
			continue
		}

		delegates[i] = username
	}

	log.Println(delegates)
	log.Println(ledger.Blocks[ledger.Users["genesis_1"]])

	/*
		for i := 0; i < 50; i++ {
			user := ledger.Users[strconv.Itoa(i+1)]
			delegate := ledger.Accounts[user]
			prev := ledger.LatestBlock(delegate.IBAN)
			blueprint, err := delegate.CreateDelegateBlock(true, prev)

			if err != nil {
				log.Println(err)
				continue
			}

			output := make(chan primitives.Block)
			node.Input <- core.NewRequest(delegate, blueprint, output)
			block := <-output
			close(output)

			if block == nil {
				log.Println("Unable to forge block")
				continue
			}
		}

		if err := dpos.Elect(account, delegates, ledger, node); err != nil {
			log.Println(err)
		}

		if err := dpos.Elect(account, []string{"-1"}, ledger, node); err != nil {
			log.Println(err)
		}
	*/
}
