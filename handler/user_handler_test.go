package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"go.test/config"
	"go.test/model"
	"go.test/repository"
	"go.test/usecase"
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
	db := config.InitDB()
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	suite.UserHandler = NewUserHandler(userUsecase)
	suite.Echo = echo.New()
}

func (suite *UserHandlerTestSuite) TestRegisterUser() {
	user := model.User{
		Username: "jae",
		Password: "123",
		Role:     "user",
	}
	body, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.Echo.NewContext(req, rec)

	err := suite.UserHandler.RegisterUser(c)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, rec.Code)
}

func (suite *UserHandlerTestSuite) TestLoginUser() {
	user := model.User{
		Username: "jae",
		Password: "123",
		Role:     "user",
	}
	db := config.InitDB()
	userRepo := repository.NewUserRepository(db)
	userRepo.Create(&user)

	loginDetails := map[string]string{
		"username": "jae",
		"password": "123",
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

func TestUpdateUser(t *testing.T) {
	e := echo.New()
	userUsecase := new(mocks.UserUsecase)
	h := NewUserHandler(userUsecase)

	user := &model.User{
		ID:       1,
		Username: "updateduser",
		Password: "newpassword",
		Role:     "admin",
	}

	jsonUser, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPut, "/restricted/users/1", bytes.NewReader(jsonUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(user.ID)))
	c.Set("role", "supervisor")

	userUsecase.On("UpdateUser", mock.Anything).Return(nil).Once()

	if assert.NoError(t, h.UpdateUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var updatedUser model.User
		json.Unmarshal(rec.Body.Bytes(), &updatedUser)
		assert.Equal(t, user.Username, updatedUser.Username)
	}

	// Test with insufficient role
	c.Set("role", "user")
	rec = httptest.NewRecorder()
	c.Response().Writer = rec

	if assert.Error(t, h.UpdateUser(c)) {
		assert.Equal(t, http.StatusForbidden, rec.Code)
	}

	// Test with higher role
	c.Set("role", "manager")
	rec = httptest.NewRecorder()
	c.Response().Writer = rec

	if assert.NoError(t, h.UpdateUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var updatedUser model.User
		json.Unmarshal(rec.Body.Bytes(), &updatedUser)
		assert.Equal(t, user.Username, updatedUser.Username)
	}

	userUsecase.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	e := echo.New()
	userUsecase := new(mocks.UserUsecase)
	h := NewUserHandler(userUsecase)

	req := httptest.NewRequest(http.MethodDelete, "/restricted/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("role", "manager")

	userUsecase.On("DeleteUser", uint(1)).Return(nil).Once()

	if assert.NoError(t, h.DeleteUser(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}

	// Test with insufficient role
	c.Set("role", "supervisor")
	rec = httptest.NewRecorder()
	c.Response().Writer = rec

	if assert.Error(t, h.DeleteUser(c)) {
		assert.Equal(t, http.StatusForbidden, rec.Code)
	}

	userUsecase.AssertExpectations(t)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
