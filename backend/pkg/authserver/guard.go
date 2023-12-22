package authserver

import (
	"fmt"
	"os/user"
	"path"

	"github.com/lestrrat-go/jwx/v2/jws"

	"github.com/julianstephens/license-server/backend/pkg/logger"
	"github.com/julianstephens/license-server/backend/pkg/model"
	"github.com/julianstephens/license-server/backend/pkg/service"
)

func Guard(conf *model.Config) error {
	// start authentication flow
	usr, err := user.Current()
	if err != nil {
		return service.HandleError(err, "failed to locate the user's home directory: %+v", nil)
	}

	// ensure token file exists & load token
	tokenDir := path.Join(usr.HomeDir, "."+conf.App.Filename)
	err = service.Ensure(tokenDir, true)
	if err != nil {
		return service.HandleError(err, "unable to ensure token dir: %+v", nil)
	}
	tokenPath := path.Join(tokenDir, ".token.json")
	err = service.Ensure(tokenPath, false)
	if err != nil {
		return service.HandleError(err, "unable to ensure token path: %+v", nil)
	}
	token, err := service.LoadToken(tokenPath)
	if err != nil {
		return service.HandleError(err, "unable to load stored token: %+v", nil)
	}

	// retrieve signing keys
	signingKeys := service.NewJWKSet(conf.Auth.JwksUrl)
	if err != nil {
		return service.HandleError(err, "unable to load signing keys for Auth0", nil)
	}

	// start auth server
	isAuthed := make(chan *string)
	authServer := *New(tokenPath, isAuthed, signingKeys, conf)
	go authServer.Start()

	// verify saved token
	var ok bool
	_, err = jws.Verify([]byte(token), jws.WithKeySet(signingKeys, jws.WithInferAlgorithmFromKey(true)))
	if err == nil {
		ok = true
		fmt.Println("existing token found & valid", token)
	} else {
		fmt.Println("existing token found & not valid", token)
		ok = false
	}

	if !ok {
		authServer.AuthenticateUser(service.NewTrue())
	} else {
		authServer.Token = &token
	}

	if authServer.Token != nil {
		logger.Infof("authenticated successfully")
	} else {
		logger.Infof("did not authenticate successfully")
	}

	return nil
}
