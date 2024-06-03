package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.test/model"
	"go.test/usecase/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	Echo        *echo.Echo
	UserHandler *UserHandler
}

func (suite *UserHandlerTestSuite) SetupSuite() {
	userUsecase := new(mocks.UserUsecase)
	suite.UserHandler = NewUserHandler(userUsecase)
	suite.Echo = echo.New()
}

func TestRegisterUser(t *testing.T) {
	e := echo.New()

	userUsecase := new(mocks.UserUsecase)
	h := NewUserHandler(userUsecase)

	user := model.User{
		Username: "ahmad",
		Password: "123",
		Role:     "user",
	}
	body, _ := json.Marshal(user)

	userUsecase.On("RegisterUser", mock.Anything).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.RegisterUser(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdUser model.User
	err = json.Unmarshal(rec.Body.Bytes(), &createdUser)

	assert.NoError(t, err)
	assert.Equal(t, user.Username, createdUser.Username)

	userUsecase.AssertExpectations(t)
}

func (suite *UserHandlerTestSuite) TestLoginUser() {

	loginDetails := map[string]string{
		"username": "rolemanager",
		"password": "jaelani",
	}
	body, _ := json.Marshal(loginDetails)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.Echo.NewContext(req, rec)

	err := suite.UserHandler.LoginUser(c)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
