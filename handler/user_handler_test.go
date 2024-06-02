package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.test/config"
	"go.test/model"
	"go.test/repository"
	"go.test/usecase"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
