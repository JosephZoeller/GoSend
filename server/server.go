package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/JosephZoeller/gmg/pkg/transit"
)

var transferAddr string

//var displayAddr string

func main() {

	transferAddr = os.Args[1]
	//displayAddr = os.Args[2]

	//go hostSave(displayAddr)

	conn, er := getSession(transferAddr)
	if er != nil {
		log.Println("Get Session Error: ", er)
		return
	}
	c := *conn
	defer c.Close()

	fHeader, er := transit.HeaderInbound(conn)
	if er != nil {
		log.Println("Create Error: ", er)
		return
	}

	er = transit.FileInbound(fHeader, conn) 
	if er != nil {
		log.Println(er)
		return
	}

	//appendSave(fHeader)
}

// Listens to a port, awaiting a connection.
func getSession(address string) (*net.Conn, error) {
	ln, er := net.Listen("tcp", address)
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
