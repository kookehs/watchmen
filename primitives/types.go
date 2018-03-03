package primitives

import (
	"crypto/sha256"
	"math/big"
)

// Amount is a type alias for big.Float that is used to represent balances for blocks.
type Amount = *big.Float

// NewAmount returns an initialized Amount with the given amount.
func NewAmount(amt float64) Amount {
	return big.NewFloat(amt)
}

// BlockHash is a sha256 hash of a block.
type BlockHash [sha256.Size]byte

// BlockHashZero is the zero value of a BlockHash
var BlockHashZero = [sha256.Size]byte{}

// BlockType is used to represent different block types in the smallest primitive possible.
type BlockType uint8

// Various block types
const (
	Change BlockType = iota
	Delegate
	Open
	Receive
	Send
)
