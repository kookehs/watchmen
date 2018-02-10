package primitives

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"io"
	"time"
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
	Timestamp() int64
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
	Time      int64           `json:"timestamp"`
}

func MakeChangeBlock(d []AccountHash, p BlockHash) ChangeBlock {
	return ChangeBlock{
		Hashables: MakeChangeHashables(d, p),
		Time:      time.Now().UnixNano(),
	}
}

func (cb *ChangeBlock) Delegates() []AccountHash {
	return cb.Hashables.Delegates
}

func (cb *ChangeBlock) Hash() BlockHash {
	var buffer bytes.Buffer
	cb.Serialize(&buffer)
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

func (cb *ChangeBlock) Timestamp() int64 {
	return cb.Time
}

func (cb *ChangeBlock) Type() BlockType {
	return Change
}

func (cb *ChangeBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(cb); err != nil {
		return err
	}

	return nil
}

func (cb *ChangeBlock) DeserializeJson(r io.Reader) error {
	return cb.Hashables.DeserializeJson(r)
}

func (cb *ChangeBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(cb); err != nil {
		return err
	}

	return nil
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

type OpenBlock struct {
	Hashables OpenHashables `json:"hashables"`
	Time      int64         `json:"timestamp"`
}

func MakeOpenBlock(a AccountHash) OpenBlock {
	return OpenBlock{
		Hashables: MakeOpenHashables(a),
		Time:      time.Now().UnixNano(),
	}
}

func (ob *OpenBlock) Delegates() []AccountHash {
	return nil
}

func (ob *OpenBlock) Hash() BlockHash {
	var buffer bytes.Buffer
	ob.Serialize(&buffer)
	return sha256.Sum256(buffer.Bytes())
}

func (ob *OpenBlock) Previous() BlockHash {
	return BlockHash{}
}

func (ob *OpenBlock) Root() BlockHash {
	return BlockHash{}
}

func (ob *OpenBlock) Source() BlockHash {
	return BlockHash{}
}

func (ob *OpenBlock) Timestamp() int64 {
	return ob.Time
}

func (ob *OpenBlock) Type() BlockType {
	return Open
}

func (ob *OpenBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(ob); err != nil {
		return err
	}

	return nil
}

func (ob *OpenBlock) DeserializeJson(r io.Reader) error {
	return ob.Hashables.DeserializeJson(r)
}

func (ob *OpenBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(ob); err != nil {
		return err
	}

	return nil
}

func (ob *OpenBlock) SerializeJson(w io.Writer) error {
	return ob.Hashables.SerializeJson(w)
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
	Time      int64            `json:"timestamp"`
}

func MakeReceiveBlock(p, s BlockHash) ReceiveBlock {
	return ReceiveBlock{
		Hashables: MakeReceiveHashables(p, s),
		Time:      time.Now().UnixNano(),
	}
}

func (rb *ReceiveBlock) Delegates() []AccountHash {
	return nil
}

func (rb *ReceiveBlock) Hash() BlockHash {
	var buffer bytes.Buffer
	rb.Serialize(&buffer)
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

func (rb *ReceiveBlock) Timestamp() int64 {
	return rb.Time
}

func (rb *ReceiveBlock) Type() BlockType {
	return Receive
}

func (rb *ReceiveBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(rb); err != nil {
		return err
	}

	return nil
}

func (rb *ReceiveBlock) DeserializeJson(r io.Reader) error {
	return rb.Hashables.DeserializeJson(r)
}

func (rb *ReceiveBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(rb); err != nil {
		return err
	}

	return nil
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
	Time      int64         `json:"timestamp"`
}

func MakeSendBlock(b Amount, d AccountHash, p BlockHash) SendBlock {
	return SendBlock{
		Hashables: MakeSendHashables(b, d, p),
		Time:      time.Now().UnixNano(),
	}
}

func (sb *SendBlock) Delegates() []AccountHash {
	return nil
}

func (sb *SendBlock) Hash() BlockHash {
	var buffer bytes.Buffer
	sb.Serialize(&buffer)
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

func (sb *SendBlock) Timestamp() int64 {
	return sb.Time
}

func (sb *SendBlock) Type() BlockType {
	return Send
}

func (sb *SendBlock) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(sb); err != nil {
		return err
	}

	return nil
}

func (sb *SendBlock) DeserializeJson(r io.Reader) error {
	return sb.Hashables.DeserializeJson(r)
}

func (sb *SendBlock) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(sb); err != nil {
		return err
	}

	return nil
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
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(ch); err != nil {
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
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(ch); err != nil {
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

type OpenHashables struct {
	Account AccountHash `json:"account"`
}

func MakeOpenHashables(a AccountHash) OpenHashables {
	return OpenHashables{
		Account: a,
	}
}

func (oh *OpenHashables) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(oh); err != nil {
		return err
	}

	return nil
}

func (oh *OpenHashables) DeserializeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(oh); err != nil {
		return err
	}

	return nil
}

func (oh *OpenHashables) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(oh); err != nil {
		return err
	}

	return nil
}

func (oh *OpenHashables) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(oh); err != nil {
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
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(rh); err != nil {
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
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(rh); err != nil {
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
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(sh); err != nil {
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
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(sh); err != nil {
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
