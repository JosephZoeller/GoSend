package main

import (
	"log"
	"net"
	"os"
)

const savefilename string = "save.json"

var server string
var port string

func main() {

	server = os.Args[1]
	port = os.Args[2]
	
	//go HostSave("8081")
	
	conn, er := getSession(port)
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
func getSession(port string) (*net.Conn, error) {
	ln, er := net.Listen("tcp", "localhost:"+port)
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
