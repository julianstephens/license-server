package authserver

import (
	"github.com/julianstephens/license-server/backend/pkg/logger"
	"github.com/julianstephens/license-server/backend/pkg/model"
	"github.com/julianstephens/license-server/backend/pkg/service"
)

func Guard(conf *model.Config) (bool, error) {
	ok := false

	_, tokenPath, err := service.GetSecureFilePath(conf, service.NewTrue())
	if err != nil {
		return false, err
	}

	token, signingKeys, err := service.GetAuthToken(conf, tokenPath)
	if err == nil {
		ok = true
	}

	// start auth server
	isAuthed := make(chan *string)
	authServer := *New(tokenPath, isAuthed, signingKeys, conf)
	go authServer.Start()

	if !ok {
		authServer.AuthenticateUser(service.NewTrue())
	} else {
		tokenStr := string(token)
		authServer.Token = &tokenStr
	}

	if authServer.Token != nil {
		logger.Infof("authenticated successfully")
		return true, nil
	} else {
		logger.Infof("did not authenticate successfully")
		return false, nil
	}
}
