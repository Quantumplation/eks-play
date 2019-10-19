package main

import (
	"log"
	"net/http"
)

func main() {
	_, err := http.Get("http://localhost:8080")
	log.Print(err)
}
