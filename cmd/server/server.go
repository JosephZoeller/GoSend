package main

import (
	"log"
	"os"

	"github.com/JosephZoeller/gmg/pkg/transit"
	"github.com/JosephZoeller/gmg/pkg/connect"
)

var transferAddr string

//var displayAddr string

func main() {

	transferAddr = os.Args[1]
	//displayAddr = os.Args[2]

	//go hostSave(displayAddr)

	conn, er := connect.OpenConnection(transferAddr)
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
