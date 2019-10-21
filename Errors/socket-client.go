package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"syscall"
)

func main() {
	_, err := http.Get("http://localhost:8080")
	if errors.Is(err, io.EOF) {
		log.Print("Error is an EOF error")
	}
	if errors.Is(err, syscall.ECONNRESET) {
		log.Print("Error is a true syscall.ECONNRESET")
	}
	if errors.Is(err, syscall.ECONNABORTED) {
		log.Print("Error is syscall.ECONNABORTED")
	}
	if errors.Is(err, syscall.ECONNREFUSED) {
		log.Print("Error is syscall.ECONNREFUSED")
	}
	log.Print(err)
}
