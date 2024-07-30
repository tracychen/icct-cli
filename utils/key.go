package utils

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

// Function to get the public key and EVM address from an ECDSA private key
func GetPublicKeyAndAddress(privateKey *ecdsa.PrivateKey) (string, string, error) {
	// Get the public key
	publicKey := privateKey.PublicKey

	// Get the uncompressed public key bytes
	publicKeyBytes := crypto.FromECDSAPub(&publicKey)

	// Compute the Keccak-256 hash of the public key (excluding the first byte)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	address := hash.Sum(nil)[12:]

	return fmt.Sprintf("%x", publicKeyBytes), fmt.Sprintf("0x%x", address), nil
}