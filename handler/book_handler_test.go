package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.test/config"
	"go.test/middleware"
	"go.test/model"
	"go.test/repository"
	"go.test/usecase"
	"go.test/usecase/mocks"
	util "go.test/utils"

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

func (suite *BookHandlerTestSuite) TestGetBooks() {
	// Generate a JWT token for testing
	token, err := util.GenerateJWT("zai", "manager")
	assert.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := suite.Echo.NewContext(req, rec)
	c.SetPath("/api/books")

	// Ensure middleware chain is properly invoked
	if assert.NoError(suite.T(), middleware.JWTMiddleware(suite.BookHandler.GetBooks)(c)) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)
	}
}

func TestCreateBook(t *testing.T) {
	e := echo.New()

	bookUsecase := new(mocks.BookUsecase)

	h := NewBookHandler(bookUsecase)

	book := &model.Book{Title: "New Book", Author: "New Author", ISBN: "0987654321", PublishedDate: "2023-01-01"}

	bookJSON, _ := json.Marshal(book)

	bookUsecase.On("CreateBook", book).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewReader(bookJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.CreateBook(c)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdBook model.Book
	err = json.Unmarshal(rec.Body.Bytes(), &createdBook)

	assert.NoError(t, err)
	assert.Equal(t, book, &createdBook)

	bookUsecase.AssertExpectations(t)
}

func TestUpdateBook(t *testing.T) {
	e := echo.New()
	bookUsecase := new(mocks.BookUsecase)
	token, err := util.GenerateJWT("zai", "supervisor")
	fmt.Println("token", token)
	h := NewBookHandler(bookUsecase)

	// Prepare a sample book for update
	book := &model.Book{
		ID:            1,
		Title:         "Updated Book",
		Author:        "Updated Author",
		ISBN:          "0987654321",
		PublishedDate: "2023-01-01",
	}

	// Convert book to JSON
	bookJSON, _ := json.Marshal(book)

	// Mock UpdateBook method to return nil error
	bookUsecase.On("UpdateBook", book).Return(nil)

	// Create a request to update a book
	req := httptest.NewRequest(http.MethodPut, "/api/books/1", bytes.NewReader(bookJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("role", "supervisor")

	// Call the UpdateBook handler
	err = h.UpdateBook(c)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the response status code is OK
	assert.Equal(t, http.StatusOK, rec.Code)
	// Decode the response body
	var updatedBook model.Book
	err = json.Unmarshal(rec.Body.Bytes(), &updatedBook)

	// Assert that the decoded book matches the input book
	assert.NoError(t, err)
	assert.Equal(t, book, &updatedBook)

	// Assert that the UpdateBook method was called once
	bookUsecase.AssertExpectations(t)
}

func TestDeleteBook(t *testing.T) {
	e := echo.New()
	bookUsecase := new(mocks.BookUsecase)
	h := NewBookHandler(bookUsecase)

	// Mock DeleteBook method to return nil error
	bookUsecase.On("DeleteBook", uint(1)).Return(nil)

	// Create a request to delete a book
	req := httptest.NewRequest(http.MethodDelete, "/api/books/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("role", "manager")

	// Call the DeleteBook handler
	err := h.DeleteBook(c)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the response status code is NoContent
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Assert that the DeleteBook method was called once
	bookUsecase.AssertExpectations(t)
}

func TestBookHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(BookHandlerTestSuite))
}
