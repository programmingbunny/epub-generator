package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	service "github.com/programmingbunny/epub-generator/service"
)

func main() {
	router := mux.NewRouter()

	//routes
	service.CreateBook(router) //add this

	log.Fatal(http.ListenAndServe(":3001", router))
}
