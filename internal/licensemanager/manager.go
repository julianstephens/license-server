package licensemanager

import (
	"crypto/elliptic"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	"github.com/julianstephens/license-server/internal/model"
	"github.com/julianstephens/license-server/internal/service"
	"github.com/julianstephens/license-server/pkg/crypto"
)

type LicenseManager struct {
	Config *model.Config
}

var ec = crypto.New(elliptic.P256())

func (lm *LicenseManager) CreateProductKeyPair(name string, productId string) (keypair *model.ProductKeyPair, err error) {
	privKey, pubKey, err := ec.GenerateKeys()
	if err != nil {
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	encodedPriv, err := ec.EncodePrivate(privKey)
	if err != nil {
		return
	}
	encodedPub, err := ec.EncodePublic(pubKey)
	if err != nil {
		return
	}

	keypair.Id = id.String()
	keypair.ProductId = productId
	keypair.PrivateKey = encodedPriv
	keypair.PublicKey = encodedPub

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

// func (lm *LicenseManager) GenerateLicense(input model.LicenseRequest) (license *model.LicenseWithAttributes, err error) {
// 	// initialize new license with user input
// 	startDate := service.If(input.StartDate != nil, *input.StartDate, time.Now().Unix())
// 	endDate := time.Unix(startDate, 10).AddDate(lm.Config.Server.LicenseLength, 0, 0).Unix()

// 	featStr := "*"
// 	if input.Features != nil {
// 		var features string
// 		features, err = jsoniter.MarshalToString(input.Features)
// 		if err != nil {
// 			return
// 		}
// 		featStr = features
// 	}

// 	var licenseKey *ecdsa.PrivateKey
// 	licenseKey, _, err = ec.GenerateKeys()
// 	if err != nil {
// 		return
// 	}

// 	licenseAttributes := fmt.Sprintf(`{"identity": "%s", "machine": "", "startDate": "%d", "endDate": "%d", "issueDate": "%d", "features": "%s"}`, input.ProductId, startDate, endDate, time.Now().Unix(), featStr)
// 	license = &model.LicenseWithAttributes{
// 		ProductId:  input.ProductId,
// 		Attributes: datatypes.JSON([]byte(licenseAttributes)),
// 	}

// 	// convert license to buffer and encrypt with product key
// 	var licenseBuf []byte
// 	licenseBuf, err = jsoniter.Marshal(&license)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var pkp *model.ProductKeyPair
// 	pkp, err = service.LoadKey(input.ProductId, lm.Config)
// 	if err != nil {
// 		err = fmt.Errorf("unable to load product key pair: %+v", err)
// 		return
// 	}

// 	// generate key pair for license issuing
// 	var priv *ecdsa.PrivateKey
// 	var pub *ecdsa.PublicKey
// 	priv, pub, err = ec.GenerateKeys()
// 	if err != nil {
// 		return
// 	}

// 	// encPriv, err := ec.EncodePrivate(priv)
// 	// if err != nil {
// 	// 	return nil, nil
// 	// }

// 	return
// }

func AssignLicense()   {}
func ValidateLicense() {}
