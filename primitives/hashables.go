package primitives

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"time"
)

type Hashables interface {
	// Deserialization
	Deserialize(io.Reader) error
	DeserializeJson(io.Reader) error

	// Serialization
	Serialize(io.Writer) error
	SerializeJson(io.Writer) error
}

type ChangeHashables struct {
	Balance   Amount        `json:"balance"`
	Delegates []AccountHash `json:"delegates"`
	Previous  BlockHash     `json:"previous"`
	Timestamp int64         `json:"timestamp"`
}

func MakeChangeHashables(a Amount, d []AccountHash, p BlockHash) ChangeHashables {
	return ChangeHashables{
		Balance:   a,
		Delegates: d,
		Previous:  p,
		Timestamp: time.Now().UnixNano(),
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
	Account   AccountHash `json:"account"`
	Balance   Amount      `json:"balance"`
	Timestamp int64       `json:"timestamp"`
}

func MakeOpenHashables(a AccountHash, b Amount) OpenHashables {
	return OpenHashables{
		Account:   a,
		Balance:   b,
		Timestamp: time.Now().UnixNano(),
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
	Balance   Amount    `json:"balance"`
	Previous  BlockHash `json:"previous"`
	Source    BlockHash `json:"source"`
	Timestamp int64     `json:"timestamp"`
}

func MakeReceiveHashables(a Amount, p, s BlockHash) ReceiveHashables {
	return ReceiveHashables{
		Balance:   a,
		Previous:  p,
		Source:    s,
		Timestamp: time.Now().UnixNano(),
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
	Timestamp   int64       `json:"timestamp"`
}

func MakeSendHashables(a Amount, d AccountHash, p BlockHash) SendHashables {
	return SendHashables{
		Balance:     a,
		Destination: d,
		Previous:    p,
		Timestamp:   time.Now().UnixNano(),
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
