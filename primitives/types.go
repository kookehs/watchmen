package primitives

import (
	"crypto/sha256"
	"math/big"
)

type Amount = big.Int

// Value of BlockHash is a sha256 hash
type BlockHash [sha256.Size]byte

var BlockHashZero = [sha256.Size]byte{}

type BlockType uint8

const (
	Change BlockType = iota
	Open
	Receive
	Send
)
