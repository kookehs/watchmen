package main

import (
	"bytes"
	"log"

	"github.com/kookehs/watchmen/core"
)

func main() {
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
