package primitives

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"time"
)

// Hashables is an interface that contains common function shared between hashable structures.
type Hashables interface {
	// Deserialization
	Deserialize(io.Reader) error
	DeserializeJSON(io.Reader) error

	// Serialization
	Serialize(io.Writer) error
	SerializeJSON(io.Writer) error
}

// ChangeHashables contains elements of a ChangeBlock that can be hashed.
type ChangeHashables struct {
	Balance   Amount    `json:"balance"`
	Delegates []IBAN    `json:"delegates"`
	Previous  BlockHash `json:"previous"`
	Timestamp int64     `json:"timestamp"`
	Type      BlockType `json:"type"`
}

// MakeChangeHashables creates and initializes a ChangeHashables from the given arguments.
func MakeChangeHashables(amt Amount, delegates []IBAN, prev BlockHash) ChangeHashables {
	return ChangeHashables{
		Balance:   amt,
		Delegates: delegates,
		Previous:  prev,
		Timestamp: time.Now().UnixNano(),
		Type:      Change,
	}
}

// Deserialize decodes byte data encoded by gob.
func (ch *ChangeHashables) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(ch)
}

// DeserializeJSON decodes JSON data.
func (ch *ChangeHashables) DeserializeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(ch)
}

// Serialize encodes to byte data using gob.
func (ch *ChangeHashables) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(ch)
}

// SerializeJSON encodes to JSON data.
func (ch *ChangeHashables) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(ch)
}

// DelegateHashables contains elements of a DelegateBlock that can be hashed.
type DelegateHashables struct {
	Balance   Amount    `json:"balance"`
	Previous  BlockHash `json:"previous"`
	Share     float64   `json:"share"`
	Timestamp int64     `json:"timestamp"`
	Type      BlockType `json:"type"`
}

// MakeDelegateHashables creates and initializes a DelegateHashables from the given arguments.
func MakeDelegateHashables(amt Amount, prev BlockHash, share float64) DelegateHashables {
	return DelegateHashables{
		Balance:   amt,
		Previous:  prev,
		Share:     share,
		Timestamp: time.Now().UnixNano(),
		Type:      Delegate,
	}
}

// Deserialize decodes byte data encoded by gob.
func (dh *DelegateHashables) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(dh)
}

// DeserializeJSON decodes JSON data.
func (dh *DelegateHashables) DeserializeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(dh)
}

// Serialize encodes to byte data using gob.
func (dh *DelegateHashables) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(dh)
}

// SerializeJSON encodes to JSON data.
func (dh *DelegateHashables) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(dh)
}

// OpenHashables contains elements of a OpenBlock that can be hashed.
type OpenHashables struct {
	Account   IBAN      `json:"account"`
	Balance   Amount    `json:"balance"`
	Timestamp int64     `json:"timestamp"`
	Type      BlockType `json:"type"`
}

// MakeOpenHashables creates and initializes a OpenHashables from the given arguments.
func MakeOpenHashables(amt Amount, iban IBAN) OpenHashables {
	return OpenHashables{
		Account:   iban,
		Balance:   amt,
		Timestamp: time.Now().UnixNano(),
		Type:      Open,
	}
}

// Deserialize decodes byte data encoded by gob.
func (oh *OpenHashables) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(oh)
}

// DeserializeJSON decodes JSON data.
func (oh *OpenHashables) DeserializeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(oh)
}

// Serialize encodes to byte data using gob.
func (oh *OpenHashables) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(oh)
}

// SerializeJSON encodes to JSON data.
func (oh *OpenHashables) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(oh)
}

// ReceiveHashables contains elements of a ReceiveBlock that can be hashed.
type ReceiveHashables struct {
	Balance   Amount    `json:"balance"`
	Previous  BlockHash `json:"previous"`
	Source    BlockHash `json:"source"`
	Timestamp int64     `json:"timestamp"`
	Type      BlockType `json:"type"`
}

// MakeReceiveHashables creates and initializes a ReceiveHashables from the given arguments.
func MakeReceiveHashables(amt Amount, prev, src BlockHash) ReceiveHashables {
	return ReceiveHashables{
		Balance:   amt,
		Previous:  prev,
		Source:    src,
		Timestamp: time.Now().UnixNano(),
		Type:      Receive,
	}
}

// Deserialize decodes byte data encoded by gob.
func (rh *ReceiveHashables) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(rh)
}

// DeserializeJSON decodes JSON data.
func (rh *ReceiveHashables) DeserializeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(rh)
}

// Serialize encodes to byte data using gob.
func (rh *ReceiveHashables) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(rh)
}

// SerializeJSON encodes to JSON data.
func (rh *ReceiveHashables) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(rh)
}

// SendHashables contains elements of a SendBlock that can be hashed.
type SendHashables struct {
	Balance     Amount    `json:"balance"`
	Destination IBAN      `json:"destination"`
	Previous    BlockHash `json:"previous"`
	Timestamp   int64     `json:"timestamp"`
	Type        BlockType `json:"type"`
}

// MakeSendHashables creates and initializes a SendHashables from the given arguments.
func MakeSendHashables(amt Amount, dst IBAN, prev BlockHash) SendHashables {
	return SendHashables{
		Balance:     amt,
		Destination: dst,
		Previous:    prev,
		Timestamp:   time.Now().UnixNano(),
		Type:        Send,
	}
}

// Deserialize decodes byte data encoded by gob.
func (sh *SendHashables) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(sh)
}

// DeserializeJSON decodes JSON data.
func (sh *SendHashables) DeserializeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(sh)
}

// Serialize encodes to byte data using gob.
func (sh *SendHashables) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(sh)
}

// SerializeJSON encodes to JSON data.
func (sh *SendHashables) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(sh)
}
