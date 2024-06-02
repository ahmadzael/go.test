package handler

import (
	"net/http"
	"strconv"

	"go.test/middleware"
	"go.test/model"
	"go.test/usecase"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase}
}

func (h *UserHandler) RegisterUser(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := h.UserUsecase.RegisterUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) LoginUser(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	token, err := h.UserUsecase.LoginUser(user.Username, user.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.UserUsecase.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.UserUsecase.GetUserByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	return middleware.RoleBasedAccess(h.updateUserHandler, "supervisor")(c)
}

func (h *UserHandler) updateUserHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	user.ID = uint(id)
	if err := h.UserUsecase.UpdateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	return middleware.RoleBasedAccess(h.deleteUserHandler, "manager")(c)
}

func (h *UserHandler) deleteUserHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.UserUsecase.DeleteUser(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
