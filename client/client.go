package main

import (
	"errors"
	"log"
	"net"
	"os"
	"time"
)

var user string = os.Getenv("USER")

func main() {

	fileIn, er := os.Open(os.Args[1])
	if er != nil {
		log.Println(er)
		return
	}

	conn, er := establishConnection(os.Args[2], 5)
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

// Anticipates a connection with the port. Times out after 30 seconds.
func establishConnection(port string, timeout int) (*net.Conn, error) {
	for i := 0; i <= timeout; i++ {
		c, er := net.Dial("tcp", "localhost:"+port)
		if er == nil {
			log.Println("[Establish Connection]: Connected")
			return &c, nil
		}
		time.Sleep(time.Second)
	}
	return nil, errors.New("[Establish Connection]: Connection Timed Out")
}
