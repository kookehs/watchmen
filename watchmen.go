package main

import (
	"crypto/elliptic"
	"encoding/gob"
	"log"
	"strconv"

	"github.com/kookehs/watchmen/core"
)

func main() {
	gob.Register(elliptic.P256())
	// TODO: Implement network for blockchain.

	// Move the below to respective test files.
	ledger := core.NewLedger()
	account, err := ledger.OpenAccount("kookehs")

	if err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < 101; i++ {
		if _, err := ledger.OpenAccount(strconv.Itoa(i + 1)); err != nil {
			log.Println(err)
			return
		}
	}

	delegates := []string{
		"+1", "+2", "+3", "+4", "+5", "+6", "+7", "+8", "+9", "+10", "+11",
		"+12", "+13", "+14", "+15", "+16", "+17", "+18", "+19", "+20", "+21", "+22",
		"+23", "+24", "+25", "+26", "+27", "+28", "+29", "+30", "+31", "+32", "+33",
		"+34", "+35", "+36", "+37", "+38", "+39", "+40", "+41", "+42", "+43", "+44",
		"+45", "+46", "+47", "+48", "+49", "+50", "+51", "+52", "+53", "+54", "+55",
		"+56", "+57", "+58", "+59", "+60", "+61", "+62", "+63", "+64", "+65", "+66",
		"+67", "+68", "+69", "+70", "+71", "+72", "+73", "+74", "+75", "+76", "+77",
		"+78", "+79", "+80", "+81", "+82", "+83", "+84", "+85", "+86", "+87", "+88",
		"+89", "+90", "+91", "+92", "+93", "+94", "+95", "+96", "+97", "+98", "+99",
		"+100", "+101",
	}

	delegate := ledger.Accounts["1"]

	prev := ledger.LatestBlock(delegate.IBAN)

	hash, err := prev.Hash()

	if err != nil {
		log.Println(err)
		return
	}

	block, err := delegate.CreateDelegateBlock(prev.Balance(), true, hash)

	if err := ledger.AppendBlock(block, delegate.IBAN); err != nil {
		log.Println(err)
		return
	}

	log.Println(ledger.Blocks[delegate.IBAN])
	core.Delegate(account, delegates, ledger)
	log.Println(account.IBAN.String())
	log.Println(account.Delegate)
	log.Println(ledger.Blocks[account.IBAN])
}
