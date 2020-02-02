package connect

import (
	"errors"
	"log"
	"net"
	"time"
)

// Listens to a port, awaiting a connection.
func OpenConnection(address string) (*net.Conn, error) {
	ln, er := net.Listen("tcp", address)
	defer ln.Close()
	if er != nil {
		log.Println(er)
		return nil, er
	}

	conn, er := ln.Accept()
	if er != nil {
		log.Println(er)
		return nil, er
	}

	return &conn, er
}

// Anticipates a connection with the port. Times out after t seconds.
func SeekConnection(addr string, t int) (*net.Conn, error) {
	for i := 0; i <= t; i++ {
		c, er := net.Dial("tcp", addr)
		if er == nil {
			log.Println("[Establish Connection]: Connected")
			return &c, nil
		}
		time.Sleep(time.Second)
	}
	return nil, errors.New("[Establish Connection]: Connection Timed Out")
}
