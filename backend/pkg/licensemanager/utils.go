package licensemanager

import (
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/x509"
	"fmt"

	"github.com/julianstephens/license-server/backend/pkg/crypto"
	"github.com/julianstephens/license-server/backend/pkg/service"
)

func (lm *LicenseManager) getKeyPair(productId string) (privKey ed25519.PrivateKey, pubKey ed25519.PublicKey, err error) {
	// retrieve product/version key pair
	pkp, err := service.LoadKey(productId, lm.Config)
	if err != nil {
		err = fmt.Errorf("unable to load product key pair: %+v", err)
		return
	}
	privKey, err = crypto.DecodePrivate(pkp.PrivateKey)
	if err != nil {
		return
	}

	pubKey, err = crypto.DecodePublic(pkp.PublicKey)

	return
}

func (lm *LicenseManager) getCipherBlock(privKey ed25519.PrivateKey) (block cipher.Block, err error) {
	genericKey, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return
	}

	key, err := crypto.Hash(string(genericKey))
	if err != nil {
		return
	}

	block, err = crypto.GetBlock(key)

	return
}
