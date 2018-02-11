package primitives

import (
	"crypto/sha256"
	"math/big"
)

// Value of AccountHash is a sha256 hash
type AccountHash [sha256.Size]byte

var AccountHashZero = [sha256.Size]byte{}

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
