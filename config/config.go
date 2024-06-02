package config

import (
	"fmt"

	"go.test/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "springstudent:springstudent@tcp(127.0.0.1:6033)/book_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database!")
		panic("Failed to connect to database!")
	}
	db.AutoMigrate(&model.Book{}, &model.User{})
	return db
}
