package main

import (
	"log"
	"net"
)

func main() {
	ln, _ := net.Listen("tcp", ":8081")
	for {
		conn, _ := ln.Accept()
		log.Print("Connection received")
		conn.Close()
		log.Print("Connection gracefully closed")
	}
}
