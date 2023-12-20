package crypto

import (
	"crypto/ed25519"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io"
)

// GenerateEd25519Key creates a secure key pair using Elliptic Curve Cryptography (ECC)
func GenerateEd25519Key() (privKey ed25519.PrivateKey, pubKey ed25519.PublicKey, err error) {
	pubKey, privKey, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return
	}

	return
}

// EncodePrivate converts a private ecdsa key to string
func EncodePrivate(privKey ed25519.PrivateKey) (key string, err error) {
	encoded, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "ED25519 PRIVATE KEY", Bytes: encoded})

	key = string(pemEncoded)

	return
}

// EncodePublic converts a public ecdsa key to string
func EncodePublic(pubKey ed25519.PublicKey) (key string, err error) {
	encoded, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return
	}

	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded})

	key = string(pemEncodedPub)
	return
}

// DecodePrivate converts a private key string to an ecdsa private key
func DecodePrivate(pemEncodedPriv string) (privKey ed25519.PrivateKey, err error) {
	blockPriv, _ := pem.Decode([]byte(pemEncodedPriv))
	x509EncodedPriv := blockPriv.Bytes
	genericPrivKey, err := x509.ParsePKCS8PrivateKey(x509EncodedPriv)
	if err != nil {
		return
	}

	privKey = genericPrivKey.(ed25519.PrivateKey)

	return
}

// DecodePublic converts a public key string to an ecdsa public key
func DecodePublic(pemEncodedPub string) (publicKey ed25519.PublicKey, err error) {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))

	x509EncodedPub := blockPub.Bytes

	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return
	}

	publicKey = genericPublicKey.(ed25519.PublicKey)

	return
}

// Sign signs the hash with the ed25519 private key
func Sign(signhash []byte, privKey ed25519.PrivateKey) (signature []byte) {
	signature = ed25519.Sign(privKey, signhash)
	return
}

// VerifySignature verifies a given hash with an ed25199 key pair
func VerifySignature(signhash []byte, privKey ed25519.PrivateKey, pubKey ed25519.PublicKey) (signature []byte, ok bool) {
	signature = Sign(signhash, privKey)
	ok = ed25519.Verify(pubKey, signhash, signature)

	return
}

// Hash computes a hash of the given data using the MD5 hash algorithm
func Hash(data string) (hash []byte, err error) {
	h := md5.New()

	_, err = io.WriteString(h, data)
	if err != nil {
		return
	}

	hash = h.Sum(nil)
	return
}
