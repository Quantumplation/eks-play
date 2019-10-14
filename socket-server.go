package main

import (
	"net"
	"os"
)

func main() {
	ln, _ := net.Listen("tcp", ":8081")
	for {
		conn, _ := ln.Accept()
		os.Exit(1)
		conn.Close()
	}
}
