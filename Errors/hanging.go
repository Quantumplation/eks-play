package main

import (
	"log"
	"net"
)

func main() {
	ln, _ := net.Listen("tcp", ":8081")
	for {
		ln.Accept()
		log.Print("Connection received, leaving it hanging")
	}
}
