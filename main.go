package main

import (
	"log"
	"net/http"

	"cli-client/controllers"
)

func main() {
	go controllers.CheckAllDevice()

	router := NewRouter()
    	log.Fatal(http.ListenAndServe(":8080", router))
}
