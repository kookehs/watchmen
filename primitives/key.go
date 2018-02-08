package primitives

import (
	"crypto/ecdsa"

	"github.com/google/uuid"
	"github.com/kookehs/watchmen/crypto"
)

type Key struct {
	Id         uuid.UUID
	Address    Address
	PrivateKey *ecdsa.PrivateKey
}

func NewKeyFromECDSA(pk *ecdsa.PrivateKey) *Key {
	return &Key{
		Id:         uuid.NewRandom(),
		Address:    crypto.ECDSAPublicKeyToSHA256(pk.PublicKey),
		PrivateKey: pk,
	}
}
