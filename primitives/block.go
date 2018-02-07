package watchmen

import (
	"bytes"
	"crypto/sha256"
	"io"
	"log"
)

// Hash values of blocks are 256 bits
type BlockHash [sha256.Size]byte

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
	Deserialize(io.Reader)
	DeserializeJson(io.Reader)

	// Serialization
	Serialize(io.Writer)
	SerializeJson(io.Writer)

	// Conversion
	String() string
	ToJson() string
}

type ChangeBlock struct {
	Hashables ChangeHashables
}

func MakeChangeBlock(d []AccountHash, p BlockHash) ChangeBlock {
	return ChangeBlock{
		Hashables: MakeChangeHashables(d, p),
	}
}

func (cb *ChangeBlock) Delegates() []AccountHash {
	return cb.Hashables.Delegates
}

func (cb *ChangeBlock) Hash() BlockHash {
	var buffer bytes.Buffer
	cb.Hashables.Serialize(&buffer)
	return sha256.Sum256(buffer.Bytes())
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

type ReceiveBlock struct {
	Hashables ReceiveHashables
}

type SendBlock struct {
	Hashables SendHashables
}

type Hashables interface {
	// Deserialization
	Deserialize(io.Reader)

	// Serialization
	Serialize(io.Writer)
}

type ChangeHashables struct {
	Delegates []AccountHash
	Previous  BlockHash
}

func MakeChangeHashables(d []AccountHash, p BlockHash) ChangeHashables {
	return ChangeHashables{
		Delegates: d,
		Previous:  p,
	}
}

func (ch *ChangeHashables) Deserialize(r io.Reader) {
	for _, d := range ch.Delegates {
		if _, err := r.Read(d[:]); err != nil {
			log.Println(err)
		}
	}

	if _, err := r.Read(ch.Previous[:]); err != nil {
		log.Println(err)
	}
}

func (ch *ChangeHashables) Serialize(w io.Writer) {
	for _, d := range ch.Delegates {
		if _, err := w.Write(d[:]); err != nil {
			log.Println(err)
		}
	}

	if _, err := w.Write(ch.Previous[:]); err != nil {
		log.Println(err)
	}
}

type ReceiveHashables struct {
	Previous BlockHash
	Source   BlockHash
}

func MakeReceiveHashables(p, s BlockHash) ReceiveHashables {
	return ReceiveHashables{
		Previous: p,
		Source:   s,
	}
}

func (rh *ReceiveHashables) Deserialize(r io.Reader) {
	if _, err := r.Read(rh.Previous[:]); err != nil {
		log.Println(err)
	}

	if _, err := r.Read(rh.Source[:]); err != nil {
		log.Println(err)
	}
}

func (rh *ReceiveHashables) Serialize(w io.Writer) {
	if _, err := w.Write(rh.Previous[:]); err != nil {
		log.Println(err)
	}

	if _, err := w.Write(rh.Source[:]); err != nil {
		log.Println(err)
	}
}

type SendHashables struct {
	Balance     Amount
	Destination AccountHash
	Previous    BlockHash
}

func MakeSendHashables(b Amount, d AccountHash, p BlockHash) SendHashables {
	return SendHashables{
		Balance:     b,
		Destination: d,
		Previous:    p,
	}
}

func (sh *SendHashables) Deserialize(r io.Reader) {
	var bytes []byte

	if _, err := r.Read(bytes); err != nil {
		log.Println(err)
	} else {
		sh.Balance.SetBytes(bytes)
	}

	if _, err := r.Read(sh.Destination[:]); err != nil {
		log.Println(err)
	}

	if _, err := r.Read(sh.Previous[:]); err != nil {
		log.Println(err)
	}
}

func (sh *SendHashables) Serialize(w io.Writer) {
	if _, err := w.Write(sh.Balance.Bytes()); err != nil {
		log.Println(err)
	}

	if _, err := w.Write(sh.Destination[:]); err != nil {
		log.Println(err)
	}

	if _, err := w.Write(sh.Previous[:]); err != nil {
		log.Println(err)
	}
}
