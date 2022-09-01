package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chapter struct {
	ChapterNum int    `json:"chapterNum" bson:"chapterNum"`
	Title      string `json:"title" bson:"title"`
	Text       string `json:"text" bson:"text"`
}

type Chapters struct {
	BookID   primitive.ObjectID `json:"bookID,omitempty" bson:"bookID,omitempty"`
	Chapters []Chapter          `json:"chapters" bson:"chapters"`
}

type ImageHeader struct {
	ImageLocation string `json:"imageLocation" bson:"imageLocation"`
	ChapterNum    int    `json:"chapterNum" bson:"chapterNum"`
}
