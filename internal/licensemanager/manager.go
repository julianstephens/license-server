package licensemanager

import (
	"bytes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/x509"
	"fmt"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/mr-tron/base58"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/julianstephens/license-server/internal/model"
	"github.com/julianstephens/license-server/internal/service"
	"github.com/julianstephens/license-server/pkg/crypto"
	"github.com/julianstephens/license-server/pkg/logger"
)

type LicenseManager struct {
	Config      *model.Config
	DB          *gorm.DB
	CurrentUser *string
}

// SetCurrentUser updates the license manager with a new gin context for user logging
func (lm *LicenseManager) SetCurrentUser(userId string) {
	lm.CurrentUser = &userId
}

// CreateProductKeyPair generates a new ecc key pair for a given product/version
func (lm *LicenseManager) CreateProductKeyPair(productId string) (keypair *model.ProductKeyPair, err error) {
	privKey, pubKey, err := crypto.GenerateEd25519Key()
	if err != nil {
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	encodedPriv, err := crypto.EncodePrivate(privKey)
	if err != nil {
		return
	}
	encodedPub, err := crypto.EncodePublic(pubKey)
	if err != nil {
		return
	}

	keypair = &model.ProductKeyPair{
		Id:         id.String(),
		ProductId:  productId,
		PrivateKey: encodedPriv,
		PublicKey:  encodedPub,
	}

	bytes, err := jsoniter.Marshal(keypair)
	if err != nil {
		return
	}

	var kpMap map[string]string
	if err = jsoniter.Unmarshal(bytes, &kpMap); err != nil {
		return
	}

	if err = service.UpdateKeyFile(kpMap, productId, false, lm.Config); err != nil {
		return
	}

	return
}

// GenerateLicense adds a new empty license to the database and returns it
func (lm *LicenseManager) GenerateLicense(productId string) (licenseBuf bytes.Buffer, eciesEd25519Enc []byte, key string, err error) {
	logger.Infof("generating license for product <%s> as user <%s>", productId, *lm.CurrentUser)

	licenseAttrs := map[string]string{
		"identity":     productId,
		"machine":      "",
		"issue_date":   "",
		"end_date":     "",
		"refresh_date": "",
	}
	attrsBuf, err := jsoniter.Marshal(licenseAttrs)
	if err != nil {
		return
	}
	attrsJson, err := jsoniter.MarshalToString(licenseAttrs)
	if err != nil {
		return
	}
	attrsHash, err := crypto.Hash(attrsJson)
	if err != nil {
		return
	}

	privKey, _, err := lm.getKeyPair(productId)
	if err != nil {
		return
	}

	sig := crypto.Sign(attrsHash, privKey)
	block, err := lm.getCipherBlock(privKey)
	if err != nil {
		return
	}

	encSig, err := crypto.Encrypt(block, string(sig))
	if err != nil {
		return
	}

	key = base58.Encode(encSig)

	license := &model.License{
		ProductId:  productId,
		Key:        []byte(key),
		Attributes: datatypes.JSON(attrsBuf),
	}
	_, err = service.Create[model.License](lm.DB, *license)

	return
}

// ValidateLicense verifies a product key has been signed by the server
func (lm *LicenseManager) ValidateKey(key string) (ok bool, err error) {
	license, err := service.Find[model.License](lm.DB, model.License{Key: []byte(key)}, nil)
	if err != nil {
		logger.Errorf("unable to retrieve license from db: %+v", err)
		return
	}
	logger.Infof("validating license for product <%s> as user <%s>", license.ProductId, *lm.CurrentUser)

	privKey, pubKey, err := lm.getKeyPair(license.ProductId)
	if err != nil {
		return
	}

	block, err := lm.getCipherBlock(privKey)
	if err != nil {
		logger.Errorf("cannot generate aes encryption key: %+v", err)
		return
	}

	decodedKey, err := base58.Decode(key)
	if err != nil {
		logger.Errorf("cannot decode product key: %+v", err)
		return
	}

	// logger.Infof("decoded product key: %+v", string(decodedKey))

	decryptedSig, err := crypto.Decrypt(block, decodedKey)
	if err != nil {
		logger.Errorf("unable to decrypt signature: %+v", err)
		return
	}

	_, ok = crypto.VerifySignature([]byte(decryptedSig), privKey, pubKey)

	return

}

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
