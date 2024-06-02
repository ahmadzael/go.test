package model

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"-"`
	Role     string `json:"role"`
	Token    string `josn:"token"`
}
