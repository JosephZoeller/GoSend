package main

import (
	"errors"
	"log"
	"net"
	"os"
	"time"
)

var user string = os.Getenv("USER")
var transferAddr string
var filename string

func main() {
	filename = os.Args[1]
	transferAddr = os.Args[2]

	fileIn, er := os.Open(filename)
	if er != nil {
		log.Println(er)
		return
	}

	conn, er := establishConnection(transferAddr, 15)
	if er != nil {
		log.Println(er)
		return
	}
	c := *conn
	defer c.Close()

	tHeader, er := makeHeader(fileIn)
	if er != nil {
		log.Println(er)
		return
	}
	er = headerOutbound(tHeader, conn)
	if er != nil {
		log.Println(er)
		return
	}

	er = fileOutbound(fileIn, conn, tHeader)
	if er != nil {
		log.Println(er)
	}
}

// Anticipates a connection with the port. Times out after t seconds.
func establishConnection(addr string, t int) (*net.Conn, error) {
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
