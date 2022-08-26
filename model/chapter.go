package model

type Chapter struct {
	ChapterNum int    `json:"chapterNum" bson:"chapterNum"`
	Title      string `json:"title" bson:"title"`
	Text       string `json:"text" bson:"text"`
}
