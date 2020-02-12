package connect

import (
	"fmt"
	//"log"
	"net"
	"time"
)

// OpenConnection Listens to an address, anticipating a connection.
func OpenConnection(address string) (*net.Conn, error) {
	ln, er := net.Listen("tcp", address)
	if er != nil {
		return nil, er
	}
	defer ln.Close()

	conn, er := ln.Accept()
	if er != nil {
		return nil, er
	}

	return &conn, er
}

// SeekConnection Attempts to connect with the address every second. Times out after t seconds.
func SeekConnection(address string, t int) (*net.Conn, error) {
	for i := 0; i <= t; i++ {
		c, er := net.Dial("tcp", address)
		if er == nil {
			return &c, nil
		}
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("Timeout after %d Seconds", t)
}
