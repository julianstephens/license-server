package controller_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/julianstephens/license-server/internal/controller"
	"github.com/julianstephens/license-server/internal/router"
	"github.com/julianstephens/license-server/pkg/database"
	"github.com/julianstephens/license-server/pkg/httputil"
	"github.com/julianstephens/license-server/pkg/model"
)

type ControllerTestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *gorm.DB
	ctrl   *controller.Controller
}

var (
	email    = "test@test.com"
	password = "test123"

	mockUser = model.AuthRequest{
		Name:     "Test User",
		Email:    email,
		Password: password,
	}

	mockTokenReq = model.AuthRequest{
		Email:    email,
		Password: password,
	}
)

func performRequest(r http.Handler, method, path string, data *bytes.Buffer) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, data)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func createUser(r http.Handler) *httptest.ResponseRecorder {
	userBuf, _ := jsoniter.Marshal(mockUser)
	w := performRequest(r, "POST", "/api/v1/auth/register", bytes.NewBuffer(userBuf))

	return w
}

func createToken(r http.Handler) *httptest.ResponseRecorder {
	reqBuf, _ := jsoniter.Marshal(mockTokenReq)
	w := performRequest(r, "POST", "/api/v1/auth/token", bytes.NewBuffer(reqBuf))

	return w
}

func issueLicense(r http.Handler, req model.LicenseRequest) *httptest.ResponseRecorder {
	reqBuf, _ := jsoniter.Marshal(req)
	w := performRequest(r, "POST", "/api/v1/licenses/issue", bytes.NewBuffer(reqBuf))

	return w
}

func handleAuth(r http.Handler) (user model.User, token model.DisplayAPIKey) {
	w := createUser(r)
	var userResp httputil.HTTPResponse[model.User]
	jsoniter.UnmarshalFromString(w.Body.String(), &userResp)
	user = userResp.Data

	w = createToken(r)

	var resp httputil.HTTPResponse[model.DisplayAPIKey]
	jsoniter.UnmarshalFromString(w.Body.String(), &resp)
	token = resp.Data

	return
}

func cleanLogs() {
	os.Remove("ls.access.log")
	os.Remove("ls.db.log")
}

func (suite *ControllerTestSuite) SetupTest() {
	suite.router = router.Setup()
	suite.db = database.GetDB()
	suite.ctrl = &controller.Controller{DB: suite.db}
	user := model.User{
		Email: email,
	}

	if err := suite.db.Where(&user).Delete(&model.User{}).Error; err != nil {
		log.Fatalf("unable to clear user: %+v", err)
	}
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
	cleanLogs()
}
