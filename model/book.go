package model

type Book struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	ISBN          string `json:"isbn"`
	PublishedDate string `json:"published_date"`
}
