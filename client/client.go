package main

import (
	"log"
	"os"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

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

	conn, er := connect.EstablishConnection(transferAddr, 15)
	if er != nil {
		log.Println(er)
		return
	}
	c := *conn
	defer c.Close()

	er = transit.HeaderOutbound(fileIn, conn)
	if er != nil {
		log.Println(er)
		return
	}

	er = transit.FileOutbound(fileIn, conn)
	if er != nil {
		log.Println(er)
	}
}
