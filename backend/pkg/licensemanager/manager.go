package licensemanager

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/mr-tron/base58"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/julianstephens/license-server/backend/pkg/crypto"
	"github.com/julianstephens/license-server/backend/pkg/logger"
	"github.com/julianstephens/license-server/backend/pkg/model"
	"github.com/julianstephens/license-server/backend/pkg/service"
)

type LicenseManager struct {
	Config *model.Config
	DB     *gorm.DB
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
func (lm *LicenseManager) GenerateLicense(productId string) (license *model.License, formattedKey string, err error) {
	logger.Infof("generating license for product <%s>", productId)

	licenseAttrs := map[string]string{
		"identity":    productId,
		"machine":     "",
		"issueDate":   "",
		"endDate":     "",
		"refreshDate": "",
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

	formattedKey = base58.Encode(encSig)

	license = &model.License{
		ProductId:  productId,
		Key:        encSig,
		Attributes: datatypes.JSON(attrsBuf),
	}
	license, err = service.Create[model.License](lm.DB, *license)

	return
}

// ValidateLicense verifies a product key has been signed by the server
func (lm *LicenseManager) ValidateKey(key string) (ok bool, license *model.License, err error) {
	decodedKey, err := base58.Decode(key)
	if err != nil {
		logger.Errorf("cannot decode product key: %+v", err)
		return
	}

	license, err = service.Find[model.License](lm.DB, model.License{Key: decodedKey}, nil)
	if err != nil {
		logger.Errorf("unable to retrieve license from db: %+v", err)
		return
	}
	logger.Infof("validating license for product <%s>", license.ProductId)

	var attrsJson map[string]interface{}
	err = jsoniter.UnmarshalFromString(license.Attributes.String(), &attrsJson)
	if err != nil {
		logger.Errorf("failed to convert attrs to go map: %+v", err)
		return
	}

	if license.Revoked || attrsJson["issueDate"] != "" {
		logger.Errorf("license for product key <%s> already issued", key)
		err = fmt.Errorf("invalid product key, already activated")
		return
	}

	privKey, pubKey, err := lm.getKeyPair(license.ProductId)
	if err != nil {
		return
	}

	block, err := lm.getCipherBlock(privKey)
	if err != nil {
		logger.Errorf("cannot generate aes encryption key: %+v", err)
		return
	}

	decryptedSig, err := crypto.Decrypt(block, decodedKey)
	if err != nil {
		logger.Errorf("unable to decrypt signature: %+v", err)
		return
	}

	_, ok = crypto.VerifySignature([]byte(decryptedSig), privKey, pubKey)

	return

}

// AssignLicense binds a license to a given machine
func (lm *LicenseManager) AssignLicense(machine string, license *model.License) (activationRes *model.ActivationData, err error) {
	logger.Infof("assigning license for product <%s>", license.ProductId)

	var attrsJson map[string]interface{}
	err = jsoniter.UnmarshalFromString(license.Attributes.String(), &attrsJson)
	if err != nil {
		logger.Errorf("failed to convert attrs to go map: %+v", err)
		return
	}

	product, err := service.FindById[model.Product](lm.DB, license.ProductId)
	if err != nil {
		logger.Errorf("unable to find product in db: %+v", err)
		return
	}

	productInfo := product.Name + " " + product.Version

	now := time.Now()
	issueDate := now.Unix()
	endDate := now.AddDate(0, 0, lm.Config.Server.DefaultLicenseLength).Unix()
	refreshDate := now.AddDate(0, 0, lm.Config.Server.MaxOfflineDuration).Unix()

	activationRes = &model.ActivationData{
		Id:          license.ExternalId,
		Product:     productInfo,
		Key:         base58.Encode(license.Key),
		IssueDate:   issueDate,
		EndDate:     endDate,
		RefreshDate: refreshDate,
	}

	if err = lm.DB.Model(&model.License{}).Where(license).UpdateColumn("attributes", datatypes.JSONSet("attributes").Set("{issueDate}", issueDate).Set("{refreshDate}", refreshDate).Set("{endDate}", endDate).Set("{machine}", machine)).Error; err != nil {
		logger.Errorf("unable to update license attributes in db: %+v", err)
		return
	}

	return
}

// RevokeLicense marks a license revoked and deletes its product key
func (lm *LicenseManager) RevokeLicense(id uuid.UUID) (err error) {
	logger.Infof("revoking license <%s>", id)

	updates := &model.License{
		Key:     []byte{},
		Revoked: true,
	}

	if err = lm.DB.Where(&model.License{ExternalId: id}).UpdateColumns(updates).Error; err != nil {
		logger.Errorf("unable to revoke license in db: %+v", err)
		return
	}

	return
}
