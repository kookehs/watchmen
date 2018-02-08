package primitives

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"io"
)

// Value of BlockHash is a sha256 hash
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
	Deserialize(io.Reader) error
	DeserializeJson(io.Reader) error

	// Serialization
	Serialize(io.Writer) error
	SerializeJson(io.Writer) error

	// Conversion
	String() (string, error)
	ToJson() (string, error)
}

type ChangeBlock struct {
	Hashables ChangeHashables `json:"hashables"`
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

func (cb *ChangeBlock) Deserialize(r io.Reader) error {
	return cb.Hashables.Deserialize(r)
}

func (cb *ChangeBlock) DeserializeJson(r io.Reader) error {
	return cb.Hashables.DeserializeJson(r)
}

func (cb *ChangeBlock) Serialize(w io.Writer) error {
	return cb.Hashables.Serialize(w)
}

func (cb *ChangeBlock) SerializeJson(w io.Writer) error {
	return cb.Hashables.SerializeJson(w)
}

func (cb *ChangeBlock) String() (string, error) {
	return cb.ToJson()
}

func (cb *ChangeBlock) ToJson() (string, error) {
	if bytes, err := json.Marshal(cb); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}

type ReceiveBlock struct {
	Hashables ReceiveHashables `json:"hashables"`
}

func MakeReceiveBlock(p, s BlockHash) ReceiveBlock {
	return ReceiveBlock{
		Hashables: MakeReceiveHashables(p, s),
	}
}

func (rb *ReceiveBlock) Delegates() []AccountHash {
	return nil
}

func (rb *ReceiveBlock) Hash() BlockHash {
	var buffer bytes.Buffer
	rb.Hashables.Serialize(&buffer)
	return sha256.Sum256(buffer.Bytes())
}

func (rb *ReceiveBlock) Previous() BlockHash {
	return rb.Hashables.Previous
}

func (rb *ReceiveBlock) Root() BlockHash {
	return rb.Hashables.Previous
}

func (rb *ReceiveBlock) Source() BlockHash {
	return rb.Hashables.Source
}

func (rb *ReceiveBlock) Type() BlockType {
	return Receive
}

func (rb *ReceiveBlock) Deserialize(r io.Reader) error {
	return rb.Hashables.Deserialize(r)
}

func (rb *ReceiveBlock) DeserializeJson(r io.Reader) error {
	return rb.Hashables.DeserializeJson(r)
}

func (rb *ReceiveBlock) Serialize(w io.Writer) error {
	return rb.Hashables.Serialize(w)
}

func (rb *ReceiveBlock) SerializeJson(w io.Writer) error {
	return rb.Hashables.SerializeJson(w)
}

func (rb *ReceiveBlock) String() (string, error) {
	return rb.ToJson()
}

func (rb *ReceiveBlock) ToJson() (string, error) {
	if bytes, err := json.Marshal(rb); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}

type SendBlock struct {
	Hashables SendHashables `json:"hashables"`
}

func MakeSendBlock(b Amount, d AccountHash, p BlockHash) SendBlock {
	return SendBlock{
		Hashables: MakeSendHashables(b, d, p),
	}
}

func (sb *SendBlock) Delegates() []AccountHash {
	return nil
}

func (sb *SendBlock) Hash() BlockHash {
	var buffer bytes.Buffer
	sb.Hashables.Serialize(&buffer)
	return sha256.Sum256(buffer.Bytes())
}

func (sb *SendBlock) Previous() BlockHash {
	return sb.Hashables.Previous
}

func (sb *SendBlock) Root() BlockHash {
	return sb.Hashables.Previous
}

func (sb *SendBlock) Source() BlockHash {
	return BlockHash{}
}

func (sb *SendBlock) Type() BlockType {
	return Send
}

func (sb *SendBlock) Deserialize(r io.Reader) error {
	return sb.Hashables.Deserialize(r)
}

func (sb *SendBlock) DeserializeJson(r io.Reader) error {
	return sb.Hashables.DeserializeJson(r)
}

func (sb *SendBlock) Serialize(w io.Writer) error {
	return sb.Hashables.Serialize(w)
}

func (sb *SendBlock) SerializeJson(w io.Writer) error {
	return sb.Hashables.SerializeJson(w)
}

func (sb *SendBlock) String() (string, error) {
	return sb.ToJson()
}

func (sb *SendBlock) ToJson() (string, error) {
	if bytes, err := json.Marshal(sb); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}

type Hashables interface {
	// Deserialization
	Deserialize(io.Reader) error
	DeserializeJson(io.Reader) error

	// Serialization
	Serialize(io.Writer) error
	SerializeJson(io.Writer) error
}

type ChangeHashables struct {
	Delegates []AccountHash `json:"delegates"`
	Previous  BlockHash     `json:"previous"`
}

func MakeChangeHashables(d []AccountHash, p BlockHash) ChangeHashables {
	return ChangeHashables{
		Delegates: d,
		Previous:  p,
	}
}

func (ch *ChangeHashables) Deserialize(r io.Reader) error {
	for _, d := range ch.Delegates {
		if _, err := r.Read(d[:]); err != nil {
			return err
		}
	}

	if _, err := r.Read(ch.Previous[:]); err != nil {
		return err
	}

	return nil
}

func (ch *ChangeHashables) DeserializeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(ch); err != nil {
		return err
	}

	return nil
}

func (ch *ChangeHashables) Serialize(w io.Writer) error {
	for _, d := range ch.Delegates {
		if _, err := w.Write(d[:]); err != nil {
			return err
		}
	}

	if _, err := w.Write(ch.Previous[:]); err != nil {
		return err
	}

	return nil
}

func (ch *ChangeHashables) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(ch); err != nil {
		return err
	}

	return nil
}

type ReceiveHashables struct {
	Previous BlockHash `json:"previous"`
	Source   BlockHash `json:"source"`
}

func MakeReceiveHashables(p, s BlockHash) ReceiveHashables {
	return ReceiveHashables{
		Previous: p,
		Source:   s,
	}
}

func (rh *ReceiveHashables) Deserialize(r io.Reader) error {
	if _, err := r.Read(rh.Previous[:]); err != nil {
		return err
	}

	if _, err := r.Read(rh.Source[:]); err != nil {
		return err
	}

	return nil
}

func (rh *ReceiveHashables) DeserializeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(rh); err != nil {
		return err
	}

	return nil
}

func (rh *ReceiveHashables) Serialize(w io.Writer) error {
	if _, err := w.Write(rh.Previous[:]); err != nil {
		return err
	}

	if _, err := w.Write(rh.Source[:]); err != nil {
		return err
	}

	return nil
}

func (rh *ReceiveHashables) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(rh); err != nil {
		return err
	}

	return nil
}

type SendHashables struct {
	Balance     Amount      `json:"balance"`
	Destination AccountHash `json:"destination"`
	Previous    BlockHash   `json:"previous"`
}

func MakeSendHashables(b Amount, d AccountHash, p BlockHash) SendHashables {
	return SendHashables{
		Balance:     b,
		Destination: d,
		Previous:    p,
	}
}

func (sh *SendHashables) Deserialize(r io.Reader) error {
	var bytes []byte

	if _, err := r.Read(bytes); err != nil {
		return err
	} else {
		sh.Balance.SetBytes(bytes)
	}

	if _, err := r.Read(sh.Destination[:]); err != nil {
		return err
	}

	if _, err := r.Read(sh.Previous[:]); err != nil {
		return err
	}

	return nil
}

func (sh *SendHashables) DeserializeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(sh); err != nil {
		return err
	}

	return nil
}

func (sh *SendHashables) Serialize(w io.Writer) error {
	if _, err := w.Write(sh.Balance.Bytes()); err != nil {
		return err
	}

	if _, err := w.Write(sh.Destination[:]); err != nil {
		return err
	}

	if _, err := w.Write(sh.Previous[:]); err != nil {
		return err
	}

	return nil
}

func (sh *SendHashables) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(sh); err != nil {
		return err
	}

	return nil
}
