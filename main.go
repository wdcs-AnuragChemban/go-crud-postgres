package main

import (
	"crud/router"
	"net/http"
	"log"
)

func main() {
	router := router.Router()

	log.Fatal(http.ListenAndServe(":8080", router))
}