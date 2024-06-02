package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	util "go.test/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MiddlewareTestSuite struct {
	suite.Suite
	Echo *echo.Echo
}

func (suite *MiddlewareTestSuite) SetupSuite() {
	suite.Echo = echo.New()
}

func (suite *MiddlewareTestSuite) TestJWTMiddleware() {
	// Valid Token
	token, _ := util.GenerateJWT("testuser", "user")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := suite.Echo.NewContext(req, rec)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}

	err := JWTMiddleware(handler)(c)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	// Invalid Token
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	rec = httptest.NewRecorder()
	c = suite.Echo.NewContext(req, rec)

	err = JWTMiddleware(handler)(c)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, rec.Code)
}

func (suite *MiddlewareTestSuite) TestRBACMiddleware() {
	token, _ := util.GenerateJWT("testadmin", "admin")
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := suite.Echo.NewContext(req, rec)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}

	// Admin role access
	err := JWTMiddleware(RBAC("admin")(handler))(c)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	// User role access denied
	token, _ = util.GenerateJWT("testuser", "user")
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	c = suite.Echo.NewContext(req, rec)

	err = JWTMiddleware(RBAC("admin")(handler))(c)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), http.StatusForbidden, rec.Code)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}
