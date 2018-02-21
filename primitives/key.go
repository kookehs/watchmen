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

// Key contains all the unique values to generate an address and private key.
type Key struct {
	ID         uuid.UUID         `json:"id"`
	Address    Address           `json:"address"`
	PrivateKey *ecdsa.PrivateKey `json:"privatekey"`
}

// NewKeyFromECDSA creates and initializes a Key generated from the given ECDSA private key.
func NewKeyFromECDSA(priv *ecdsa.PrivateKey) (*Key, error) {
	id, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	hash := crypto.ECDSAPublicKeyToSHA256(priv.PublicKey)
	address := MakeAddress(hash[:])

	return &Key{
		ID:         id,
		Address:    address,
		PrivateKey: priv,
	}, nil
}

// NewKeyForICAP creates and initializes a Key for the Inter-exchange Client Address Protocol.
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

// Deserialize decodes byte data encoded by gob.
func (k *Key) Deserialize(r io.Reader) error {
	decoder := gob.NewDecoder(r)
	return decoder.Decode(k)
}

// DeseralizeJSON decodes JSON data.
func (k *Key) DeseralizeJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(k)
}

// Serialize encodes to byte data using gob.
func (k *Key) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(k)
}

// SerializeJSON encodes to JSON data.
func (k *Key) SerializeJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(k)
}
