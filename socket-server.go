package main

import (
	"net"
	"os"
)

func main() {
	ln, _ := net.Listen("tcp", ":8081")
	for {
		_, _ = ln.Accept()
		os.Exit(1)
	}
}
