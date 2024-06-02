package config

import (
	"fmt"
	"os"

	"go.test/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database!")
		panic("Failed to connect to database!")
	}
	db.AutoMigrate(&model.Book{}, &model.User{})
	return db
}
