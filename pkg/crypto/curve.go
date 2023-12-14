package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io"
)

type EllipticCurve struct {
	pubKeyCurve elliptic.Curve // http://golang.org/pkg/crypto/elliptic/#P256
	privateKey  *ecdsa.PrivateKey
	publicKey   *ecdsa.PublicKey
}

func New(curve elliptic.Curve) *EllipticCurve {
	return &EllipticCurve{
		pubKeyCurve: curve,
		privateKey:  new(ecdsa.PrivateKey),
	}
}

// GenerateKeys creates a secure key pair using Elliptic Curve Cryptography (ECC)
func (ec *EllipticCurve) GenerateKeys() (privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, err error) {
	privateKey, err = ecdsa.GenerateKey(ec.pubKeyCurve, rand.Reader)
	if err != nil {
		return
	}

	ec.privateKey = privateKey
	ec.publicKey = &privateKey.PublicKey
	return
}

// EncodePrivate converts a private ecdsa key to string
func (ec *EllipticCurve) EncodePrivate(privKey *ecdsa.PrivateKey) (key string, err error) {
	encoded, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: encoded})

	key = string(pemEncoded)

	return
}

// EncodePublic converts a public ecdsa key to string
func (ec *EllipticCurve) EncodePublic(pubKey *ecdsa.PublicKey) (key string, err error) {
	encoded, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return
	}

	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded})

	key = string(pemEncodedPub)
	return
}

// DecodePrivate converts a private key string to an ecdsa private key
func (ec *EllipticCurve) DecodePrivate(pemEncodedPriv string) (privateKey *ecdsa.PrivateKey, err error) {
	blockPriv, _ := pem.Decode([]byte(pemEncodedPriv))
	x509EncodedPriv := blockPriv.Bytes
	privateKey, err = x509.ParseECPrivateKey(x509EncodedPriv)
	if err != nil {
		return
	}

	return
}

// DecodePublic converts a public key string to an ecdsa public key
func (ec *EllipticCurve) DecodePublic(pemEncodedPub string) (publicKey *ecdsa.PublicKey, err error) {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))

	x509EncodedPub := blockPub.Bytes

	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return
	}

	publicKey = genericPublicKey.(*ecdsa.PublicKey)

	return
}

// VerifySignature verifies a given hash with an ecdsa key pair
func (ec *EllipticCurve) VerifySignature(signhash []byte, privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) (signature []byte, ok bool, err error) {
	signature, err = ecdsa.SignASN1(rand.Reader, privKey, signhash)
	if err != nil {
		return
	}

	ok = ecdsa.VerifyASN1(pubKey, signhash, signature)

	return
}

// Hash computes a hash of the given data using the MD5 hash algorithm
func (ec *EllipticCurve) Hash(data string) (hash []byte, err error) {
	h := md5.New()

	_, err = io.WriteString(h, data)
	if err != nil {
		return
	}

	hash = h.Sum(nil)
	return
}
