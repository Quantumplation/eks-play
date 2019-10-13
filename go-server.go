package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func socketServer() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Print(err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			return
		}
		conn.Close()
	}
}

func main() {
	go socketServer()

	url := os.Getenv("URL")
	if url == "" {
		url = "http://localhost:8081/sample"
	}

	http.HandleFunc("/sample", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Request")
		fmt.Fprintf(w, "Good")
	})
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	for {
		log.Printf("Requesting %s", url)
		resp, err := http.Get(url)
		if err != nil {
			// ECONNRESET
			if errors.Unwrap(err) == io.EOF {
				log.Print("ECONNRESET!")
				log.Print(err)
			}
			return
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			return
		}
		resp.Body.Close()
	}
}
