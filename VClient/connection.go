/*
 * Author: Igor joaquim dos Santos Lima
 * Github: https://github.com/igor036
*/
package main

import (
  "net"
  "log"
)

const PROTOCOL = "tcp"

func Connection(addr string) net.Conn {

	conn, err := net.Dial(PROTOCOL, addr)

	if err != nil { log.Fatal(err) }

	return conn

}
