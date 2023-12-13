package service

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/julianstephens/license-server/internal/model"
)

type KeyFile struct {
	KeyPairs []map[string]string
}

func SaveKeyPair(data map[string]string, productId string, shouldRemove bool, conf *model.Config) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	saveLoc := path + conf.Server.KeyPath

	// load any existing product keys
	existingKeys, err := os.ReadFile(saveLoc)
	if err != nil {
		return err
	}

	var f KeyFile
	if err := jsoniter.Unmarshal(existingKeys, &f); err != nil {
		return err
	}

	for i, kp := range f.KeyPairs {
		if kp["id"] == productId {
			if shouldRemove {
				f.KeyPairs = DeleteElement(f.KeyPairs, i)
			} else {
				f.KeyPairs[i] = data
			}
		} else {
			f.KeyPairs = append(f.KeyPairs, data)
		}
	}

	byteData, err := jsoniter.Marshal(f)
	if err != nil {
		return err
	}

	if err := os.WriteFile(saveLoc, byteData, os.ModePerm); err != nil {
		return err
	}

	return nil
}
