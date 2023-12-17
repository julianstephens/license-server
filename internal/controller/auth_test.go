package controller_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

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

type ControllerTestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *gorm.DB
}

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

func cleanLogs() {
	os.Remove("ls.access.log")
	os.Remove("ls.db.log")
}

func (suite *ControllerTestSuite) SetupTest() {
	suite.router = router.Setup()
	suite.db = database.GetDB()
	user := model.User{
		Email: email,
	}

	if err := suite.db.Where(&user).Delete(&model.User{}).Error; err != nil {
		log.Fatalf("unable to clear user: %+v", err)
	}
}

func (suite *ControllerTestSuite) TestRegister() {
	w := createUser(suite.router)

	suite.Assert().Equal(http.StatusCreated, w.Code)

	var resp httputil.HTTPResponse[model.User]
	jsoniter.UnmarshalFromString(w.Body.String(), &resp)
	gotUser := resp.Data

	suite.Assert().NotZero(gotUser.ID)
	suite.Equal(mockUser.Name, gotUser.Name)
	suite.Equal(mockUser.Email, gotUser.Email)
}

func (suite *ControllerTestSuite) TestCreateToken() {
	w := createUser(suite.router)
	var userResp httputil.HTTPResponse[model.User]
	jsoniter.UnmarshalFromString(w.Body.String(), &userResp)
	user := userResp.Data

	w = createToken(suite.router)

	suite.Assert().Equal(http.StatusCreated, w.Code)

	var resp httputil.HTTPResponse[model.DisplayAPIKey]
	jsoniter.UnmarshalFromString(w.Body.String(), &resp)
	gotKey := resp.Data

	expiresAtTime, _ := time.Parse(time.DateTime, gotKey.ExpiresAt)
	expiresAt := expiresAtTime.Unix()

	suite.Assert().NotZero(gotKey.ID)
	suite.Assert().Equal(user.ID, gotKey.UserId)
	suite.Assert().Greater(expiresAt, int64(gotKey.CreatedAt))
	suite.Assert().Len(gotKey.Key, 32)
	suite.Assert().Zero(gotKey.Scopes)
}

func (suite *ControllerTestSuite) TestAuthorize() {
	w := createUser(suite.router)
	var userResp httputil.HTTPResponse[model.User]
	jsoniter.UnmarshalFromString(w.Body.String(), &userResp)
	user := userResp.Data

	w = createToken(suite.router)

	var resp httputil.HTTPResponse[model.DisplayAPIKey]
	jsoniter.UnmarshalFromString(w.Body.String(), &resp)
	token := resp.Data

	ctrl := controller.Controller{DB: suite.db}
	isAdmin, userId, _ := ctrl.Authorize(token.Key, "admin")

	suite.Assert().False(isAdmin)
	suite.Assert().Zero(userId)

	if err := suite.db.Model(&model.APIKey{}).Where(&model.APIKey{UserId: user.ID}).UpdateColumn("scopes", "admin").Error; err != nil {
		log.Fatal(err)
	}

	isAdmin, userId, _ = ctrl.Authorize(token.Key, "admin")

	suite.Assert().True(isAdmin)
	suite.Assert().Equal(user.ID, userId)
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
	cleanLogs()
}
