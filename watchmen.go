package main

import (
	"crypto/rand"
	"log"

	"github.com/kookehs/watchmen/primitives"
)

func main() {
	account, err := primitives.NewAccount(rand.Reader)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(account.Key.Address.String())
	log.Println(account.Bban.String())
	log.Println(account.Iban.String())
}
