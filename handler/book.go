package handler

import (
	"net/http"
	"strconv"

	"go.test/model"
	"go.test/usecase"

	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	BookUsecase usecase.BookUsecase
}

func NewBookHandler(bookUsecase usecase.BookUsecase) *BookHandler {
	return &BookHandler{bookUsecase}
}

func (h *BookHandler) GetBooks(c echo.Context) error {
	books, err := h.BookUsecase.GetAllBooks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetBook(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	book, err := h.BookUsecase.GetBookByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, book)
}

func (h *BookHandler) CreateBook(c echo.Context) error {
	book := new(model.Book)
	if err := c.Bind(book); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := h.BookUsecase.CreateBook(book); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	book := new(model.Book)
	if err := c.Bind(book); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	book.ID = uint(id)
	if err := h.BookUsecase.UpdateBook(book); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.BookUsecase.DeleteBook(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
