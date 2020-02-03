package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/jsonUtil"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

type save struct {
	Servers    []id `json:"Servers"`
	ProxyPorts []id `json:"ProxyPorts"`
}

type id struct {
	IP   string `json:"IP"`
	Port string `json:"Port"`
}

var saveFile = "addresses.json"
var addresses = save{}
var lastSpeak = ""

// 1.) open the listeners listed in the json
// 2.) each listener will wait for a connection
// 3.) after a connection, round robin the server ports
// 4.) when a connection ends, the listener is available for another connection

func main() {
	jsonUtil.LoadFromFile(saveFile, &addresses)
	log.Printf("Registered %d ports", len(addresses.ProxyPorts))
	for i, v := range addresses.ProxyPorts {
		go OpenListener(v.IP + v.Port)
		log.Printf("Proxy %d listening at %s", i, (v.IP + v.Port))
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	<-signalChan
}

func OpenListener(address string) {

	ln, _ := net.Listen("tcp", address)
	for {
		lCon, _ := ln.Accept()
		DoConnect(&lCon, pickSpeak())
		lCon.Close()
	}
}

func DoConnect(lCon *net.Conn, sAddr string) error {

	sCon, er := connect.SeekConnection(sAddr, 5)
	log.Println("The reverse proxy has connected to " + sAddr)
	if er != nil {
		log.Println(er)
		return er
	}

	fHead, er := transit.PassHeader(lCon, sCon)
	if er != nil {
		log.Println(er)
		return er
	}

	return transit.PassFile(fHead, lCon, sCon)
}

func pickSpeak() string {
	addrsCnt := len(addresses.Servers)

	if lastSpeak == "" {
		this := addresses.Servers[0]
		lastSpeak = this.IP + this.Port
		log.Println("Choosing " + lastSpeak + " to serve")
		return lastSpeak
	}

	for i := 0; i < addrsCnt; i++ {
		this := addresses.Servers[i]
		if lastSpeak == this.IP+this.Port {
			next := addresses.Servers[(i+1)%addrsCnt]
			lastSpeak = next.IP + next.Port
			log.Println("Choosing " + lastSpeak + " to serve")
			break
		}
	}

	return lastSpeak
}
