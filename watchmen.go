package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"log"

	"github.com/kookehs/watchmen/core"
	"github.com/kookehs/watchmen/primitives"
)

func main() {
	gob.Register(elliptic.P256())

	key, err := primitives.NewKeyForICAP(rand.Reader)

	if err != nil {
		log.Println(err)
		return
	}

	account := core.MakeAccount(*key)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(account.Key.Address.String())
	log.Println(account.BBAN.String())
	log.Println(account.IBAN.String())

	ledger := core.NewLedger()
	ledger.OpenAccount("kookehs")
	// log.Println(ledger.Accounts)
}
