package connect

import (
	"errors"
	"log"
	"net"
	"time"
)

// Listens to an address, anticipating a connection.
func OpenConnection(address string) (*net.Conn, error) {
	ln, er := net.Listen("tcp", address)
	defer ln.Close()
	if er != nil {
		log.Println(er)
		return nil, er
	}

	conn, er := ln.Accept()
	log.Println("[Open Connection]: Connection Established")
	if er != nil {
		log.Println(er)
		return nil, er
	}

	return &conn, er
}

// Attempts to connect with the address every second. Times out after t seconds.
func SeekConnection(address string, t int) (*net.Conn, error) {
	for i := 0; i <= t; i++ {
		c, er := net.Dial("tcp", address)
		if er == nil {
			log.Println("[Seek Connection]: Connection Established")
			return &c, nil
		}
		time.Sleep(time.Second)
	}
	return nil, errors.New("[Establish Connection]: Connection Timed Out")
}
