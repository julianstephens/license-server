package service

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"

	"github.com/julianstephens/license-server/backend/pkg/logger"
	"github.com/julianstephens/license-server/backend/pkg/model"
)

type KeyFile struct {
	KeyPairs *[]map[string]string
	Token    *string
}

type KeyFileOpts struct {
	IsKeypair    *bool
	IsToken      *bool
	ShouldRemove *bool
	FileLoc      *string
}

func LoadKeyPair(productId string, conf *model.Config) (*model.ProductKeyPair, error) {
	_, saveFile, err := GetSecureFilePath(conf, nil)
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
	if f.KeyPairs != nil {
		for _, val := range *f.KeyPairs {
			if val["product_id"] == productId {
				if err := mapstructure.Decode(val, &pkp); err != nil {
					return nil, err
				}
			}
		}

		return &pkp, nil

	} else {
		return nil, fmt.Errorf("no keys to load")
	}
}

func LoadToken(path string) (string, error) {
	contents, err := ReadFile(path)
	if err != nil {
		return *new(string), err
	}

	if len(contents) == 0 {
		return *new(string), HandleError(fmt.Errorf("token file empty"), "could not load token", nil)
	}

	var f KeyFile
	if err := jsoniter.Unmarshal(contents, &f); err != nil {
		logger.Errorf("could not unmarshal file to struct: %+v", err)
		return *new(string), err
	}

	return *f.Token, nil
}

func ClearToken(path string) error {
	if err := os.Truncate(path, 0); err != nil {
		return HandleError(err, "failed to truncate token file", nil)
	}

	return nil
}

// UpdateKeyFile adds, replaces, or removes a key pair from the given key file
func UpdateKeyFile[T map[string]string | string](data T, resourceID string, conf *model.Config, opts *KeyFileOpts) error {
	didUpdate := false
	var saveDir, saveFile string
	var err error

	if *opts.IsToken {
		saveDir, saveFile, err = GetSecureFilePath(conf, NewTrue())
		if err != nil {
			return err
		}
	} else {
		saveDir, saveFile, err = GetSecureFilePath(conf, NewFalse())
		if err != nil {
			return err
		}
	}

	if opts.FileLoc != nil {
		splitPath := strings.Split(*opts.FileLoc, string(os.PathSeparator))
		saveDir = path.Join(splitPath[:len(splitPath)-1]...)
		saveFile = path.Join(string(os.PathSeparator), saveDir, splitPath[len(splitPath)-1])
	}

	if err := Ensure(saveDir, true); err != nil {
		logger.Errorf("could not ensure file dir")
		return err
	}
	if err := Ensure(saveFile, false); err != nil {
		logger.Errorf("could not ensure file")
		return err
	}

	jsonFile, err := os.Open(saveFile)
	if err != nil {
		logger.Errorf("could not open file")
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

	if opts.IsKeypair != nil && *opts.IsKeypair {
		var writeData map[string]string
		err = mapstructure.Decode(data, &writeData)
		if err != nil {
			return err
		}

		keyPairs := *f.KeyPairs
		for i, kp := range keyPairs {
			if kp["product_id"] == resourceID {
				if *opts.ShouldRemove {
					keyPairs = DeleteElement(keyPairs, i)
					didUpdate = true
				} else {
					keyPairs[i] = writeData
					didUpdate = true
				}
			}
		}

		if !didUpdate && !*opts.ShouldRemove {
			keyPairs = append(keyPairs, writeData)
			f.KeyPairs = &keyPairs
		}
	} else if opts.IsToken != nil && *opts.IsToken {
		writeToken := reflect.ValueOf(data).String()

		if opts.ShouldRemove != nil && *opts.ShouldRemove {
			f.Token = nil
			didUpdate = true
		} else {
			f.Token = &writeToken
			didUpdate = true
		}
	}

	if !didUpdate && (opts.ShouldRemove != nil && !*opts.ShouldRemove) {
		msg := fmt.Sprintf("no %s found for %s: %s", If(*opts.IsKeypair || !*opts.IsToken, "key pair", "token"), If(*opts.IsKeypair || !*opts.IsToken, "product", "resource"), resourceID)
		return HandleError(fmt.Errorf(msg), msg, nil)
	}

	byteData, err := jsoniter.Marshal(f)
	if err != nil {
		return HandleError(err, "unable to convert key file to bytes", nil)
	}

	truncatedFile, err := os.Create(saveFile)
	if err != nil {
		return HandleError(err, "unable to truncate key file", nil)
	}
	defer truncatedFile.Close()

	if _, err := truncatedFile.Write(byteData); err != nil {
		logger.Errorf("unable to write to key file: %v", err)
		return err
	}

	return nil
}

func GetTokenPath(conf *model.Config) (string, error) {
	// start authentication flow
	usr, err := user.Current()
	if err != nil {
		return *new(string), HandleError(err, "failed to locate the user's home directory: %+v", nil)
	}

	// ensure token file exists & load token
	tokenDir := path.Join(usr.HomeDir, "."+conf.App.Filename)
	err = Ensure(tokenDir, true)
	if err != nil {
		return *new(string), HandleError(err, "unable to ensure token dir: %+v", nil)
	}
	tokenPath := path.Join(tokenDir, ".token.json")
	err = Ensure(tokenPath, false)
	if err != nil {
		return *new(string), HandleError(err, "unable to ensure token path: %+v", nil)
	}

	return tokenPath, nil
}

func GetSecureFilePath(conf *model.Config, isToken *bool) (keyFileDirPath string, keyFilePath string, err error) {
	p, err := os.UserHomeDir()
	if err != nil {
		return
	}
	keyFileDirPath = path.Join(p, conf.Server.OutDir)
	keyFilePath = path.Join(keyFileDirPath, If(*isToken, conf.Server.TokenFile, conf.Server.KeyFile))
	return
}
