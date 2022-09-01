package chapter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/programmingbunny/epub-generator/model"
)

// calls the epub-backend client using a book's unique _id
// this call retrieves every chapter of the book
func GetChapters(id string) model.Chapters {
	response, err := http.Get("http://localhost:3000/chapters/" + id)
	if err != nil {
		fmt.Println(err.Error())
		return model.Chapters{}
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return model.Chapters{}
	}

	var newChapter model.Chapters
	json.Unmarshal(responseData, &newChapter)

	return newChapter
}

// calls epub-backend client
// uses book's unique _id to retrieve book details (title, author name, etc.)
func GetBookDetails(id string) model.Book {
	response, err := http.Get("http://localhost:3000/book/" + id)
	if err != nil {
		fmt.Println(err.Error())
		return model.Book{}
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return model.Book{}
	}

	var newBook model.Book
	json.Unmarshal(responseData, &newBook)

	return newBook
}

// calls epub-backend client
// uses a chapter's unique _id to retrieve that chapter's data
func getChapter(id string) model.Chapter {
	response, err := http.Get("http://localhost:3000/chapter/" + id)
	if err != nil {
		fmt.Println(err.Error())
		return model.Chapter{}
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return model.Chapter{}
	}

	var newChapter model.Chapter
	json.Unmarshal(responseData, &newChapter)

	return newChapter
}

func CreateNewChapter(chapter model.Chapter) string {
	returnThis := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<!DOCTYPE html>
	<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="en" lang="en">
		<head>
			<title>` + chapter.Title + `</title>
		</head>
		<body>
				<p>` + chapter.Text + `</p>
		</body>
	</html>`
	return returnThis
}

// calls epub-backend client
// uses book's unique _id & chapterNum to retrieve chapter header's image location
func GetChapterHeaderImage(id string, num string) model.ImageHeader {
	response, err := http.Get("http://localhost:3000/getChapterImage/" + id + "/" + num)
	if err != nil {
		fmt.Println(err.Error())
		return model.ImageHeader{}
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return model.ImageHeader{}
	}

	var imageLoc model.ImageHeader
	json.Unmarshal(responseData, &imageLoc)

	return imageLoc
}
