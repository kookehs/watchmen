package watchmen

import (
	"bufio"
)

// Hash values of blocks are 256 bits
type BlockHash [32]byte

type BlockType uint8

const (
	Change BlockType = iota
	Open
	Receive
	Send
)

type Block interface {
	// Block
	Delegates() []AccountHash
	Hash() BlockHash
	Previous() BlockHash
	Root() BlockHash
	Source() BlockHash
	Type() BlockType

	// Deserialization
	Deserialize(*bufio.Reader)
	DeserializeJson(*bufio.Reader)

	// Serialization
	Serialize(*bufio.Writer)
	SerializeJson(*bufio.Writer)

	// Conversion
	String() string
	ToJson() string
}

type ChangeBlock struct {
	Hashables ChangeHashables
}

func MakeChangeBlock(previous BlockHash, delegates []AccountHash) ChangeBlock {
	return ChangeBlock{
		Hashables: MakeChangeHashables(previous, delegates),
	}
}

func (cb *ChangeBlock) Delegates() []AccountHash {
	return cb.Hashables.Delegates
}

func (cb *ChangeBlock) Hash() BlockHash {
	// TODO: Hash hashables
	return BlockHash{}
}

func (cb *ChangeBlock) Previous() BlockHash {
	return cb.Hashables.Previous
}

func (cb *ChangeBlock) Root() BlockHash {
	return cb.Hashables.Previous
}

func (cb *ChangeBlock) Source() BlockHash {
	return BlockHash{}
}

func (cb *ChangeBlock) Type() BlockType {
	return Change
}

type ChangeHashables struct {
	Previous  BlockHash
	Delegates []AccountHash
}

func MakeChangeHashables(previous BlockHash, delegates []AccountHash) ChangeHashables {
	return ChangeHashables{
		Previous:  previous,
		Delegates: delegates,
	}
}

type ReceiveBlock struct {
	Hashables ReceiveHashables
}

type ReceiveHashables struct {
	Previous BlockHash
	Source   BlockHash
}

type SendBlock struct {
	Hashables SendHashables
}

type SendHashables struct {
	Previous    BlockHash
	Destination AccountHash
	Balance     Amount
}
