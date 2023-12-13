package service

import (
	"fmt"
	"os"
	"path"

	jsoniter "github.com/json-iterator/go"
	"github.com/julianstephens/license-server/internal/model"
	"github.com/julianstephens/license-server/pkg/logger"
)

type KeyFile struct {
	KeyPairs []map[string]string
}

// UpdateKeyPairFile adds, replaces, or removes a key pair from the server's key file
func UpdateKeyPairFile(data map[string]string, productId string, shouldRemove bool, conf *model.Config) error {
	didUpdate := false

	p, err := os.Getwd()
	if err != nil {
		return err
	}

	saveDir := path.Join(p, conf.Server.OutDir)
	if err := Ensure(saveDir, true); err != nil {
		return err
	}
	saveFile := path.Join(saveDir, conf.Server.KeyFile)
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
