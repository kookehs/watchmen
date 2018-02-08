package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
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
