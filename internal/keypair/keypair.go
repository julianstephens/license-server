package keypair

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
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

func (ec *EllipticCurve) GenerateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(ec.pubKeyCurve, rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to generate key pair: %s", err)
	}

	ec.privateKey = privateKey
	ec.publicKey = &privateKey.PublicKey
	return privateKey, &privateKey.PublicKey, nil
}

func (ec *EllipticCurve) EncodePrivate(privKey *ecdsa.PrivateKey) (string, error) {
	encoded, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return "", err
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: encoded})

	key := string(pemEncoded)

	return key, nil
}

func (ec *EllipticCurve) EncodePublic(pubKey *ecdsa.PublicKey) (string, error) {
	encoded, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", nil
	}

	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded})

	key := string(pemEncodedPub)
	return key, nil
}

func (ec *EllipticCurve) DecodePrivate(pemEncodedPriv string) (*ecdsa.PrivateKey, error) {
	blockPriv, _ := pem.Decode([]byte(pemEncodedPriv))
	x509EncodedPriv := blockPriv.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509EncodedPriv)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (ec *EllipticCurve) DecodePublic(pemEncodedPub string) (*ecdsa.PublicKey, error) {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))

	x509EncodedPub := blockPub.Bytes

	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return nil, err
	}

	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey, nil
}

func (ec *EllipticCurve) VerifySignature(privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) ([]byte, bool, error) {
	h := md5.New()

	_, err := io.WriteString(h, "This is a message to be signed and verified by ECDSA!")
	if err != nil {
		return nil, false, err
	}
	signhash := h.Sum(nil)

	r, s, serr := ecdsa.Sign(rand.Reader, privKey, signhash)
	if serr != nil {
		return []byte(""), false, serr
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	ok := ecdsa.Verify(pubKey, signhash, r, s)

	return signature, ok, nil
}
