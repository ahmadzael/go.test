package usecase

import (
	"go.test/model"
	"go.test/repository"
)

type BookUsecase interface {
	GetAllBooks() ([]model.Book, error)
	GetBookByID(id uint) (*model.Book, error)
	CreateBook(book *model.Book) error
	UpdateBook(book *model.Book) error
	DeleteBook(id uint) error
}

type bookUsecase struct {
	bookRepo repository.BookRepository
}

func NewBookUsecase(bookRepo repository.BookRepository) BookUsecase {
	return &bookUsecase{bookRepo}
}

func (u *bookUsecase) GetAllBooks() ([]model.Book, error) {
	return u.bookRepo.GetAll()
}

func (u *bookUsecase) GetBookByID(id uint) (*model.Book, error) {
	return u.bookRepo.GetByID(id)
}

func (u *bookUsecase) CreateBook(book *model.Book) error {
	return u.bookRepo.Create(book)
}

func (u *bookUsecase) UpdateBook(book *model.Book) error {
	return u.bookRepo.Update(book)
}

func (u *bookUsecase) DeleteBook(id uint) error {
	return u.bookRepo.Delete(id)
}
