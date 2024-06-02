package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"go.test/config"
	"go.test/middleware"
	"go.test/model"
	"go.test/repository"
	"go.test/usecase"
	"go.test/usecase/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	c.Set("role", "user")

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
func TestUpdateBook(t *testing.T) {
	e := echo.New()
	bookUsecase := new(mocks.BookUsecase)
	h := NewBookHandler(bookUsecase)

	book := &model.Book{
		ID:            1,
		Title:         "Updated Sample Book",
		Author:        "Updated Sample Author",
		ISBN:          "0987654321",
		PublishedDate: "2023-01-01",
	}

	jsonBook, _ := json.Marshal(book)

	req := httptest.NewRequest(http.MethodPut, "/restricted/books/1", bytes.NewReader(jsonBook))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(book.ID)))
	c.Set("role", "supervisor")

	bookUsecase.On("UpdateBook", mock.Anything).Return(nil).Once()

	if assert.NoError(t, h.UpdateBook(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var updatedBook model.Book
		json.Unmarshal(rec.Body.Bytes(), &updatedBook)
		assert.Equal(t, book.Title, updatedBook.Title)
	}

	// Test with insufficient role
	c.Set("role", "user")
	rec = httptest.NewRecorder()
	c.Response().Writer = rec

	if assert.Error(t, h.UpdateBook(c)) {
		assert.Equal(t, http.StatusForbidden, rec.Code)
	}

	// Test with higher role
	c.Set("role", "manager")
	rec = httptest.NewRecorder()
	c.Response().Writer = rec

	if assert.NoError(t, h.UpdateBook(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var updatedBook model.Book
		json.Unmarshal(rec.Body.Bytes(), &updatedBook)
		assert.Equal(t, book.Title, updatedBook.Title)
	}

	bookUsecase.AssertExpectations(t)
}

func TestDeleteBook(t *testing.T) {
	e := echo.New()
	bookUsecase := new(mocks.BookUsecase)
	h := NewBookHandler(bookUsecase)

	req := httptest.NewRequest(http.MethodDelete, "/restricted/books/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("role", "manager")

	bookUsecase.On("DeleteBook", uint(1)).Return(nil).Once()

	if assert.NoError(t, h.DeleteBook(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}

	// Test with insufficient role
	c.Set("role", "supervisor")
	rec = httptest.NewRecorder()
	c.Response().Writer = rec

	if assert.Error(t, h.DeleteBook(c)) {
		assert.Equal(t, http.StatusForbidden, rec.Code)
	}

	bookUsecase.AssertExpectations(t)
}

func TestBookHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(BookHandlerTestSuite))
}
