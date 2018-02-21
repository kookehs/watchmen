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
	log.Println(account.Bban.String())
	log.Println(account.Iban.String())

	ledger := core.NewLedger()
	ledger.CreateAccount("kookehs")
	// log.Println(ledger.Accounts)
}
