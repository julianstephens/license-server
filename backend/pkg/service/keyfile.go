package service

import (
	"fmt"
	"os"
	"path"

	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"

	"github.com/julianstephens/license-server/backend/pkg/logger"
	"github.com/julianstephens/license-server/backend/pkg/model"
)

type KeyFile struct {
	KeyPairs []map[string]string
}

func LoadKey(productId string, conf *model.Config) (*model.ProductKeyPair, error) {
	_, saveFile, err := getKeyFilePaths(conf)
	if err != nil {
		return nil, err
	}

	contents, err := ReadFile(saveFile)
	if err != nil {
		return nil, err
	}

	var pkp model.ProductKeyPair
	var f KeyFile
	if err := jsoniter.Unmarshal(contents, &f); err != nil {
		return nil, err
	}

	for _, val := range f.KeyPairs {
		if val["product_id"] == productId {
			if err := mapstructure.Decode(val, &pkp); err != nil {
				return nil, err
			}
		}
	}

	return &pkp, nil
}

// UpdateKeyFile adds, replaces, or removes a key pair from the server's key file
func UpdateKeyFile(data map[string]string, productId string, shouldRemove bool, conf *model.Config) error {
	didUpdate := false

	saveDir, saveFile, err := getKeyFilePaths(conf)
	if err != nil {
		return err
	}

	if err := Ensure(saveDir, true); err != nil {
		return err
	}
	if err := Ensure(saveFile, false); err != nil {
		return err
	}

	jsonFile, err := os.Open(saveFile)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	existingKeys, err := ReadFile(saveFile)
	if err != nil {
		return err
	}

	var f KeyFile
	if len(existingKeys) > 0 {
		if err := jsoniter.Unmarshal(existingKeys, &f); err != nil {
			return err
		}
	}

	for i, kp := range f.KeyPairs {
		if kp["product_id"] == productId {
			if shouldRemove {
				f.KeyPairs = DeleteElement(f.KeyPairs, i)
				didUpdate = true
			} else {
				f.KeyPairs[i] = data
				didUpdate = true
			}
		}
	}

	if !didUpdate && !shouldRemove {
		f.KeyPairs = append(f.KeyPairs, data)
	}

	if !didUpdate && shouldRemove {
		return fmt.Errorf("no key pairs found for product: %s", productId)
	}

	byteData, err := jsoniter.Marshal(f)
	if err != nil {
		logger.Errorf("unable to convert key file to bytes: %v", err)
		return err
	}

	truncatedFile, err := os.Create(saveFile)
	if err != nil {
		logger.Errorf("unable to truncate key pair file: %v", err)
		return err
	}
	defer truncatedFile.Close()

	if _, err := truncatedFile.Write(byteData); err != nil {
		logger.Errorf("unable to write to key pair file: %v", err)
		return err
	}

	return nil
}

func getKeyFilePaths(conf *model.Config) (keyFileDirPath string, keyFilePath string, err error) {
	p, err := os.Getwd()
	keyFileDirPath = path.Join(p, conf.Server.OutDir)
	keyFilePath = path.Join(keyFileDirPath, conf.Server.KeyFile)
	return
}
