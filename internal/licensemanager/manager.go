package licensemanager

import (
	"crypto/elliptic"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/julianstephens/license-server/internal/keypair"
	"github.com/julianstephens/license-server/internal/model"
	"github.com/julianstephens/license-server/service"
)

var ec = keypair.New(elliptic.P256())

func CreateProductKeyPair(name string, conf *model.Config) (*model.ProductKeyPair, error) {
	var keypair model.ProductKeyPair

	privKey, pubKey, err := ec.GenerateKeys()
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	encodedPriv, err := ec.EncodePrivate(privKey)
	if err != nil {
		return nil, err
	}
	encodedPub, err := ec.EncodePublic(pubKey)
	if err != nil {
		return nil, err
	}

	keypair.Id = id.String()
	keypair.Name = name
	keypair.PrivateKey = encodedPriv
	keypair.PublicKey = encodedPub

	bytes, err := jsoniter.Marshal(keypair)
	if err != nil {
		return nil, err
	}

	var kpMap map[string]string
	if err := jsoniter.Unmarshal(bytes, &kpMap); err != nil {
		return nil, err
	}

	if err := service.SaveKeyPair(kpMap, keypair.Id, false, conf); err != nil {
		return nil, err
	}

	return &keypair, nil
}

func GenerateLicense() {}
func ExportLicense()   {}
func ValidateLicense() {}
