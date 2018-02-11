package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

func ECDSAPublicKeyToOctet(pk *ecdsa.PublicKey) []byte {
	if pk == nil || pk.X == nil || pk.Y == nil {
		return nil
	}

	return elliptic.Marshal(elliptic.P256(), pk.X, pk.Y)
}

func ECDSAPublicKeyToSHA256(pk ecdsa.PublicKey) [sha256.Size]byte {
	octet := ECDSAPublicKeyToOctet(&pk)
	return sha256.Sum256(octet[1:])
}

func Sign(hash []byte, pk *ecdsa.PrivateKey) (*big.Int, *big.Int, error) {
	return ecdsa.Sign(rand.Reader, pk, hash)
}

func Verify(hash []byte, pk *ecdsa.PublicKey, r, s *big.Int) bool {
	return ecdsa.Verify(pk, hash, r, s)
}
