package model

type Book struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Author    string `json:"author"`
	BookCover string `json:"bookCover"`
}
