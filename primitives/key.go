package primitives

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"io"

	"github.com/google/uuid"
	"github.com/kookehs/watchmen/crypto"
)

type Key struct {
	Id         uuid.UUID
	Address    *Address
	PrivateKey *ecdsa.PrivateKey
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
