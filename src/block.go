package watchmen

import (
	"bufio"
)

type Account = [32]byte

// Hash values of blocks are 256 bits
type BlockHash = [32]byte

type BlockType = uint8

const (
	Change BlockType = iota
	Open
	Receive
	Send
)

type Block interface {
	Delegate() []Account
	Hash() BlockHash
	PreviousHash() BlockHash
	Root() BlockHash
	Serialize(*bufio.Writer)
	SerializeJson(*bufio.Writer)
	ToJson() string
	Type() BlockType
}
