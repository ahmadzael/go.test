package main

import (
	"go.test/config"
	"go.test/handler"
	"go.test/middleware"
	"go.test/repository"
	"go.test/usecase"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	db := config.InitDB()

	bookRepo := repository.NewBookRepository(db)
	bookUsecase := usecase.NewBookUsecase(bookRepo)
	bookHandler := handler.NewBookHandler(bookUsecase)

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	e.POST("/register", userHandler.RegisterUser)
	e.POST("/login", userHandler.LoginUser)

	// JWT protected routes
	r := e.Group("/restricted")
	r.Use(middleware.JWTMiddleware)

	// Book routes
	r.GET("/books", bookHandler.GetBooks)
	r.GET("/books/:id", bookHandler.GetBook)
	r.POST("/books", bookHandler.CreateBook, middleware.RBAC("admin"))
	r.PUT("/books/:id", bookHandler.UpdateBook, middleware.RBAC("admin"))
	r.DELETE("/books/:id", bookHandler.DeleteBook, middleware.RBAC("admin"))

	// User routes
	r.GET("/users", userHandler.GetUsers, middleware.RBAC("admin"))
	r.GET("/users/:id", userHandler.GetUser, middleware.RBAC("admin"))
	r.PUT("/users/:id", userHandler.UpdateUser, middleware.RBAC("admin"))
	r.DELETE("/users/:id", userHandler.DeleteUser, middleware.RBAC("admin"))

	e.Logger.Fatal(e.Start(":1323"))
}
