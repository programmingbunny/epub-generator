package routes

import (
	"github.com/gorilla/mux"
	epub "github.com/programmingbunny/epub-generator/controller/make-epub"
)

func CreateBook(router *mux.Router) {
	router.HandleFunc("/book/{bookId}", epub.CreateBook()).Methods("GET")
}
