package usecase

import (
	"errors"

	"go.test/model"
	"go.test/repository"
	util "go.test/utils"
)

type UserUsecase interface {
	RegisterUser(user *model.User) error
	LoginUser(username, password string) (string, error)
	GetAllUsers() ([]model.User, error)
	GetUserByID(id uint) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) RegisterUser(user *model.User) error {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return u.userRepo.Create(user)
}

func (u *userUsecase) LoginUser(username, password string) (string, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return "", err
	}
	if !util.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid username or password")
	}
	token, err := util.GenerateJWT(user.Username, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userUsecase) GetAllUsers() ([]model.User, error) {
	return u.userRepo.GetAll()
}

func (u *userUsecase) GetUserByID(id uint) (*model.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *userUsecase) UpdateUser(user *model.User) error {
	return u.userRepo.Update(user)
}

func (u *userUsecase) DeleteUser(id uint) error {
	return u.userRepo.Delete(id)
}
