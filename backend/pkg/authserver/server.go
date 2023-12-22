package authserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"

	"github.com/julianstephens/license-server/backend/pkg/logger"
	"github.com/julianstephens/license-server/backend/pkg/model"
	"github.com/julianstephens/license-server/backend/pkg/service"
)

type AuthServer struct {
	ClientID     string
	ServiceAppID string
	TenantID     string
	Finished     chan *string
	Token        *string
	TokenPath    string
	SigningKeys  jwk.Set
	Conf         *model.Config
	httpServer   *http.Server
}

func New(tokenPath string, finished chan *string, signingKeys jwk.Set, conf *model.Config) *AuthServer {
	return &AuthServer{
		ClientID:     conf.Auth.ClientID,
		ServiceAppID: conf.Auth.ResourceID,
		TenantID:     conf.Auth.TenantID,
		Finished:     finished,
		TokenPath:    tokenPath,
		SigningKeys:  signingKeys,
		Conf:         conf,
	}
}

func (s *AuthServer) Start() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("could not read working dir")
	}
	staticDir := path.Join(dir, "backend", "static")
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fs)
	http.HandleFunc("/token", s.handleToken)
	logger.Infof("Auth Server listening on port %d", s.Conf.Auth.Port)
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", s.Conf.Auth.Port),
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatalf("could not start auth server: %+v", s.httpServer.ListenAndServe())
}

func (s *AuthServer) AuthenticateUser(openBrowser *bool) {
	if *openBrowser {
		var cmd *exec.Cmd
		cmd = nil
		r := runtime.GOOS
		url := fmt.Sprintf("http://localhost:%d/", s.Conf.Auth.Port)
		if strings.HasPrefix(r, "darwin") {
			// Do OSX things
			cmd = exec.Command("open", url)
		} else if strings.HasPrefix(r, "linux") {
			// Do linux things
			cmd = exec.Command("xdg-open", url)
		} else if strings.HasPrefix(r, "windows") {
			cmd = exec.Command("start", url)
		}
		if cmd == nil {
			*openBrowser = false
		} else {
			cmd.Stdin = strings.NewReader("")
			f, err := os.Open(os.DevNull)
			if err != nil {
				log.Fatalf("Error opening /dev/null?!\t%s\n", err)
				*openBrowser = false
			} else {
				cmd.Stdout = f
				err = cmd.Start()
				if err != nil {
					log.Fatalf("Unable to run the command: %s", err)
					*openBrowser = false
				}
			}
		}
	}
	if !*openBrowser {
		fmt.Printf("In order to complete authentication, open your browser and navigate to http://%s:%d/\n", s.Conf.Auth.Host, s.Conf.Auth.Port)
	}

	// block on the channel, waiting for authentication to occur
	s.Token = <-s.Finished

	logger.Infof("shutting down HTTP server")
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Printf("Unable to shutdown http server: %s", err)
	}
}

func (s *AuthServer) handleToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	_, err := jws.Verify([]byte(token), jws.WithKeySet(s.SigningKeys, jws.WithInferAlgorithmFromKey(true)))
	if err != nil {
		logger.Errorf("unable to parse token: %+v", err)
	}
	if err == nil {
		err = service.UpdateKeyFile(token, s.ServiceAppID, s.Conf, &service.KeyFileOpts{IsToken: service.NewTrue(), FileLoc: &s.TokenPath})
		if err != nil {
			log.Fatalf("failed to save token to file: %+v", err)
		} else {
			w.WriteHeader(204)
			_, err = w.Write([]byte{})
			if err != nil {
				log.Fatal(err)
			}
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			fmt.Println(token)
			time.Sleep(1000 * time.Millisecond)
			s.Finished <- &token
		}
	}
}
