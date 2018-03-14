package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"math/big"
)

// Register various types to allow encoding.
func init() {
	gob.Register(elliptic.P256())
}

// ECDSAPublicKeyToOctet returns the byte of a point in octect representation.
func ECDSAPublicKeyToOctet(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}

	return elliptic.Marshal(elliptic.P256(), pub.X, pub.Y)
}

// ECDSAPublicKeyToSHA256 returns the SHA256 hash minus the first byte.
func ECDSAPublicKeyToSHA256(pub ecdsa.PublicKey) [sha256.Size]byte {
	octet := ECDSAPublicKeyToOctet(&pub)
	return sha256.Sum256(octet[1:])
}

// Sign signs the given hash with the given private key.
func Sign(hash []byte, priv *ecdsa.PrivateKey) (*big.Int, *big.Int, error) {
	return ecdsa.Sign(rand.Reader, priv, hash)
}

// Verify verifies the given hash with the given public key.
func Verify(hash []byte, pub *ecdsa.PublicKey, r, s *big.Int) bool {
	return ecdsa.Verify(pub, hash, r, s)
}
