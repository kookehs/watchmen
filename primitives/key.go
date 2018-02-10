package primitives

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"encoding/json"
	"io"

	"github.com/google/uuid"
	"github.com/kookehs/watchmen/crypto"
)

type Key struct {
	Id         uuid.UUID         `json:"id"`
	Address    *Address          `json:"address"`
	PrivateKey *ecdsa.PrivateKey `json:"privatekey"`
}

func NewKeyFromECDSA(pk *ecdsa.PrivateKey) (*Key, error) {
	id, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	hash := crypto.ECDSAPublicKeyToSHA256(pk.PublicKey)
	address := NewAddress(hash[:])

	return &Key{
		Id:         id,
		Address:    address,
		PrivateKey: pk,
	}, nil
}

func NewKeyForICAP(r io.Reader) (*Key, error) {
	noise := make([]byte, 64)

	if _, err := r.Read(noise); err != nil {
		return nil, err
	}

	reader := bytes.NewReader(noise)
	privateKeyECDSA, err := ecdsa.GenerateKey(elliptic.P256(), reader)

	if err != nil {
		return nil, err
	}

	key, err := NewKeyFromECDSA(privateKeyECDSA)

	if err != nil {
		return nil, err
	}

	return key, nil
}

func (k *Key) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	if err := decoder.Decode(k); err != nil {
		return err
	}

	return nil
}

func (k *Key) DeseralizeJson(r io.Reader) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(k); err != nil {
		return err
	}

	return nil
}

func (k *Key) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	if err := encoder.Encode(k); err != nil {
		return err
	}

	return nil
}

func (k *Key) SerializeJson(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(k); err != nil {
		return err
	}

	return nil
}
