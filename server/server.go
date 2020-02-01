package main

import (
	"log"
	"net"
	"os"
)

const savefilename string = "save.json"

var serverNum string
var transferAddr string
var displayAddr string

func main() {

	serverNum = os.Args[1]
	transferAddr = os.Args[2]
	displayAddr = os.Args[3]
	
	go hostSave(displayAddr)
	
	conn, er := getSession(transferAddr)
	if er != nil {
		log.Println("Get Session Error: ", er)
		return
	}
	c := *conn
	defer c.Close()

	tHeader, er := headerInbound(conn)
	if er != nil {
		log.Println("Create Error: ", er)
		return
	}

	fileCreate, er := os.Create(tHeader.Filename)
	if er != nil {
		log.Println("Create Error: ", er)
		return
	}
	defer fileCreate.Close()

	if fileInbound(tHeader, conn, fileCreate) != nil {
		return
	}
	appendSave(tHeader)
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
