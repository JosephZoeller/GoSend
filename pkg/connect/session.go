package connect

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/JosephZoeller/gmg/pkg/logUtil"
)

// Listens to an address, anticipating a connection.
func OpenConnection(address string) (*net.Conn, error) {
	ln, er := net.Listen("tcp", address)
	if er != nil {
		return nil, logUtil.FormatError("Connect OpenConnection", er)
	}
	defer ln.Close()

	conn, er := ln.Accept()
	log.Println("[Connect OpenConnection]: Connection Established")
	if er != nil {
		return nil, logUtil.FormatError("Connect OpenConnection", er)
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
	return nil, logUtil.FormatError("Connect SeekConnection", errors.New(fmt.Sprintf("Timeout after %d Seconds", t)))
}
