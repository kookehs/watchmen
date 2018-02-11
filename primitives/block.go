package primitives

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"io"

	"github.com/kookehs/watchmen/crypto"
)

type Block interface {
	// Block
	Balance() Amount
	Delegates() []AccountHash
	Hash() BlockHash
	Previous() BlockHash
	Root() BlockHash
	Sign(*ecdsa.PrivateKey) error
	Source() BlockHash
	Timestamp() int64
	Type() BlockType
	Verify(*ecdsa.PublicKey) bool

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
	Signature Signature       `json:"signature"`
}

func MakeChangeBlock(a Amount, d []AccountHash, p BlockHash) ChangeBlock {
	return ChangeBlock{
		Hashables: MakeChangeHashables(a, d, p),
	}
}

func (cb *ChangeBlock) Balance() Amount {
	return cb.Hashables.Balance
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
	return BlockHashZero
}

func (cb *ChangeBlock) Sign(pk *ecdsa.PrivateKey) error {
	hash := cb.Hash()
	r, s, err := crypto.Sign(hash[:], pk)

	if err != nil {
		return err
	}

	cb.Signature = MakeSignature(r, s)
	return nil
}

func (cb *ChangeBlock) Timestamp() int64 {
	return cb.Hashables.Timestamp
}

func (cb *ChangeBlock) Type() BlockType {
	return Change
}

func (cb *ChangeBlock) Verify(pk *ecdsa.PublicKey) bool {
	hash := cb.Hash()
	return crypto.Verify(hash[:], pk, cb.Signature.R, cb.Signature.S)
}

func (cb *ChangeBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(cb); err != nil {
		return err
	}

	return nil
}

func (cb *ChangeBlock) DeserializeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(cb); err != nil {
		return err
	}

	return nil
}

func (cb *ChangeBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(cb); err != nil {
		return err
	}

	return nil
}

func (cb *ChangeBlock) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(cb); err != nil {
		return err
	}

	return nil
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

type OpenBlock struct {
	Hashables OpenHashables `json:"hashables"`
	Signature Signature     `json:"signature"`
}

func MakeOpenBlock(a AccountHash, b Amount) OpenBlock {
	return OpenBlock{
		Hashables: MakeOpenHashables(a, b),
	}
}

func (ob *OpenBlock) Balance() Amount {
	return ob.Hashables.Balance
}

func (ob *OpenBlock) Delegates() []AccountHash {
	return nil
}

func (ob *OpenBlock) Hash() BlockHash {
	var buffer bytes.Buffer
	ob.Hashables.Serialize(&buffer)
	return sha256.Sum256(buffer.Bytes())
}

func (ob *OpenBlock) Previous() BlockHash {
	return BlockHashZero
}

func (ob *OpenBlock) Root() BlockHash {
	return BlockHashZero
}

func (ob *OpenBlock) Sign(pk *ecdsa.PrivateKey) error {
	hash := ob.Hash()
	r, s, err := crypto.Sign(hash[:], pk)

	if err != nil {
		return err
	}

	ob.Signature = MakeSignature(r, s)
	return nil
}

func (ob *OpenBlock) Source() BlockHash {
	return BlockHashZero
}

func (ob *OpenBlock) Timestamp() int64 {
	return ob.Hashables.Timestamp
}

func (ob *OpenBlock) Type() BlockType {
	return Open
}

func (ob *OpenBlock) Verify(pk *ecdsa.PublicKey) bool {
	hash := ob.Hash()
	return crypto.Verify(hash[:], pk, ob.Signature.R, ob.Signature.S)
}

func (ob *OpenBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(ob); err != nil {
		return err
	}

	return nil
}

func (ob *OpenBlock) DeserializeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(ob); err != nil {
		return err
	}

	return nil
}

func (ob *OpenBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(ob); err != nil {
		return err
	}

	return nil
}

func (ob *OpenBlock) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(ob); err != nil {
		return err
	}

	return nil
}

func (ob *OpenBlock) String() (string, error) {
	return ob.ToJson()
}

func (ob *OpenBlock) ToJson() (string, error) {
	if bytes, err := json.Marshal(ob); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}

type ReceiveBlock struct {
	Hashables ReceiveHashables `json:"hashables"`
	Signature Signature        `json:"signature"`
}

func MakeReceiveBlock(a Amount, p, s BlockHash) ReceiveBlock {
	return ReceiveBlock{
		Hashables: MakeReceiveHashables(a, p, s),
	}
}

func (rb *ReceiveBlock) Balance() Amount {
	return rb.Hashables.Balance
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

func (rb *ReceiveBlock) Sign(pk *ecdsa.PrivateKey) error {
	hash := rb.Hash()
	r, s, err := crypto.Sign(hash[:], pk)

	if err != nil {
		return err
	}

	rb.Signature = MakeSignature(r, s)
	return nil
}

func (rb *ReceiveBlock) Source() BlockHash {
	return rb.Hashables.Source
}

func (rb *ReceiveBlock) Timestamp() int64 {
	return rb.Hashables.Timestamp
}

func (rb *ReceiveBlock) Type() BlockType {
	return Receive
}

func (rb *ReceiveBlock) Verify(pk *ecdsa.PublicKey) bool {
	hash := rb.Hash()
	return crypto.Verify(hash[:], pk, rb.Signature.R, rb.Signature.S)
}

func (rb *ReceiveBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(rb); err != nil {
		return err
	}

	return nil
}

func (rb *ReceiveBlock) DeserializeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(rb); err != nil {
		return err
	}

	return nil
}

func (rb *ReceiveBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(rb); err != nil {
		return err
	}

	return nil
}

func (rb *ReceiveBlock) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(rb); err != nil {
		return err
	}

	return nil
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
	Signature Signature     `json:"signature"`
}

func MakeSendBlock(a Amount, d AccountHash, p BlockHash) SendBlock {
	return SendBlock{
		Hashables: MakeSendHashables(a, d, p),
	}
}

func (sb *SendBlock) Balance() Amount {
	return sb.Hashables.Balance
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

func (sb *SendBlock) Sign(pk *ecdsa.PrivateKey) error {
	hash := sb.Hash()
	r, s, err := crypto.Sign(hash[:], pk)

	if err != nil {
		return err
	}

	sb.Signature = MakeSignature(r, s)
	return nil
}

func (sb *SendBlock) Source() BlockHash {
	return BlockHashZero
}

func (sb *SendBlock) Timestamp() int64 {
	return sb.Hashables.Timestamp
}

func (sb *SendBlock) Type() BlockType {
	return Send
}

func (sb *SendBlock) Verify(pk *ecdsa.PublicKey) bool {
	hash := sb.Hash()
	return crypto.Verify(hash[:], pk, sb.Signature.R, sb.Signature.S)
}

func (sb *SendBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(sb); err != nil {
		return err
	}

	return nil
}

func (sb *SendBlock) DeserializeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(sb); err != nil {
		return err
	}

	return nil
}

func (sb *SendBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(sb); err != nil {
		return err
	}

	return nil
}

func (sb *SendBlock) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(sb); err != nil {
		return err
	}

	return nil
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
