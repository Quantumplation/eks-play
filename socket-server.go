package main

import (
	"net"
)

func main() {
	ln, _ := net.Listen("tcp", ":8081")
	for {
		conn, _ := ln.Accept()
		conn.Close()
	}
}
