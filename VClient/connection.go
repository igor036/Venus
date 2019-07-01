// Author: Igor joaquim dos Santos Lima
// Github: https://github.com/igor036

package main

import (
	"log"
	"net"
)

const protocol = "tcp"

// Connection return the connection with server
func Connection(addr string) net.Conn {

	conn, err := net.Dial(protocol, addr)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	return conn

}
