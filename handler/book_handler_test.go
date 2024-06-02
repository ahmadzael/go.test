package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.test/config"
	"go.test/middleware"
	"go.test/model"
	"go.test/repository"
	"go.test/usecase"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BookHandlerTestSuite struct {
	suite.Suite
	Echo        *echo.Echo
	BookHandler *BookHandler
}

func (suite *BookHandlerTestSuite) SetupSuite() {
	db := config.InitDB()
	bookRepo := repository.NewBookRepository(db)
	bookUsecase := usecase.NewBookUsecase(bookRepo)
	suite.BookHandler = NewBookHandler(bookUsecase)
	suite.Echo = echo.New()
}

func (suite *BookHandlerTestSuite) TestCreateBook() {
	book := model.Book{
		Title:         "Test Book",
		Author:        "Test Author",
		ISBN:          "1234567890",
		PublishedDate: "2022-01-01",
	}
	body, _ := json.Marshal(book)

	req := httptest.NewRequest(http.MethodPost, "/restricted/books", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.Echo.NewContext(req, rec)
	c.Set("role", "admin")

	err := middleware.JWTMiddleware(suite.BookHandler.CreateBook)(c)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, rec.Code)

	var createdBook model.Book
	err = json.Unmarshal(rec.Body.Bytes(), &createdBook)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), book.Title, createdBook.Title)
}

func (suite *BookHandlerTestSuite) TestGetBooks() {
	req := httptest.NewRequest(http.MethodGet, "/restricted/books", nil)
	rec := httptest.NewRecorder()
	c := suite.Echo.NewContext(req, rec)
	c.Set("role", "user")

	err := middleware.JWTMiddleware(suite.BookHandler.GetBooks)(c)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
}

func TestBookHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(BookHandlerTestSuite))
}
