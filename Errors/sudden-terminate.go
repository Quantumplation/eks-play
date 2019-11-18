package main

import (
	"log"
	"net"
	"os"
)

func main() {
	ln, _ := net.Listen("tcp", ":8080")
	for {
		ln.Accept()
		log.Print("Connection received, exiting the process")
		os.Exit(1)
	}
}
