package main

import (
	"log"
	"os"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

var listenAddr string
var speakAddr string

func main() {
	listenAddr = os.Args[1]
	speakAddr = os.Args[2]

	lCon, er := connect.OpenConnection(listenAddr)
	if er != nil {
		log.Println(er)
		return
	}

	sCon, er := connect.SeekConnection(speakAddr, 30)
	if er != nil {
		log.Println(er)
		return
	}

	fHead, er := transit.PassHeader(lCon, sCon)
	if er != nil {
		log.Println(er)
		return
	}

	er = transit.PassFile(fHead, lCon, sCon)
	if er != nil {
		log.Println(er)
		return
	}

}