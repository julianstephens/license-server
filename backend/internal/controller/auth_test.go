package controller_test

import (
	"log"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/julianstephens/license-server/backend/pkg/httputil"
	"github.com/julianstephens/license-server/backend/pkg/model"
)

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
	user, token := handleAuth(suite.router)

	isAdmin, userId, _ := suite.ctrl.Authorize(token.Key, "admin")

	suite.Assert().False(isAdmin)
	suite.Assert().Zero(userId)

	if err := suite.db.Model(&model.APIKey{}).Where(&model.APIKey{UserId: user.ID}).UpdateColumn("scopes", "admin").Error; err != nil {
		log.Fatal(err)
	}

	isAdmin, userId, _ = suite.ctrl.Authorize(token.Key, "admin")

	suite.Assert().True(isAdmin)
	suite.Assert().Equal(user.ID, userId)
}
