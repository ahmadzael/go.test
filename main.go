package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"go.test/config"
	"go.test/handler"
	"go.test/middleware"
	"go.test/repository"
	"go.test/usecase"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	if env := godotenv.Load(); env != nil {
		fmt.Println(".env only for env local")
	}

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

	e.POST("/api/register", userHandler.RegisterUser)
	e.POST("/api/login", userHandler.LoginUser)

	restricted := e.Group("/api")
	restricted.Use(middleware.JWTMiddleware)

	restricted.GET("/books", bookHandler.GetBooks)
	restricted.GET("/books/:id", bookHandler.GetBook)
	restricted.POST("/books", middleware.RoleBasedAccess(bookHandler.CreateBook, "supervisor"))
	restricted.PUT("/books/:id", middleware.RoleBasedAccess(bookHandler.UpdateBook, "supervisor"))
	restricted.DELETE("/books/:id", middleware.RoleBasedAccess(bookHandler.DeleteBook, "manager"))

	restricted.GET("/users", userHandler.GetUsers)
	restricted.GET("/users/:id", userHandler.GetUser)
	restricted.PUT("/users/:id", middleware.RoleBasedAccess(userHandler.UpdateUser, "supervisor"))
	restricted.DELETE("/users/:id", middleware.RoleBasedAccess(userHandler.DeleteUser, "manager"))

	// Start server
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "1323"
	}
	s := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	e.Logger.Fatal(e.StartServer(s))
}
