package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"io"

	"github.com/kookehs/watchmen/primitives"
)

type Account struct {
	Bban primitives.BBAN `json:"bban"`
	Iban primitives.IBAN `json:"iban"`
	Key  primitives.Key  `josn:"key"`
}

func MakeAccount(k primitives.Key) Account {
	bban := primitives.MakeBBAN([]byte(k.Address.String()))
	iban := primitives.MakeIBAN([]byte("TV00" + bban.String()))

	return Account{
		Bban: bban,
		Iban: iban,
		Key:  k,
	}
}

func (a *Account) CreateChangeBlock(b primitives.Amount, d []primitives.IBAN, p primitives.BlockHash) (*primitives.ChangeBlock, error) {
	block := primitives.NewChangeBlock(b, d, p)

	if err := block.Sign(a.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

func (a *Account) CreateOpenBlock() (*primitives.OpenBlock, error) {
	var balance primitives.Amount
	block := primitives.NewOpenBlock(a.Iban, balance)

	if err := block.Sign(a.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

func (a *Account) CreateReceiveBlock(b primitives.Amount, p, s primitives.BlockHash) (*primitives.ReceiveBlock, error) {
	block := primitives.NewReceiveBlock(b, p, s)

	if err := block.Sign(a.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

func (a *Account) CreateSendBlock(b primitives.Amount, d primitives.IBAN, p primitives.BlockHash) (*primitives.SendBlock, error) {
	block := primitives.NewSendBlock(b, d, p)

	if err := block.Sign(a.Key.PrivateKey); err != nil {
		return nil, err
	}

	return block, nil
}

func (a *Account) Hash() [sha256.Size]byte {
	var buffer bytes.Buffer
	a.Serialize(&buffer)
	return sha256.Sum256(buffer.Bytes())
}

func (a *Account) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(a); err != nil {
		return err
	}

	return nil
}

func (a *Account) DeseralizeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(a); err != nil {
		return err
	}

	return nil
}

func (a *Account) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(a); err != nil {
		return err
	}

	return nil
}

func (a *Account) SeralizeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(a); err != nil {
		return err
	}

	return nil
}
