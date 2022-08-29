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
