package main

import (
	"errors"
	"log"
	"net"
	"os"
	"sync"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/jsonUtil"
	"github.com/JosephZoeller/gmg/pkg/transit"
)


type save struct {
	Servers []server `json:"Servers"`
}

type server struct {
	IP   string `json:"IP"`
	Port string `json:"Port"`
}

var serverFile = "servers.json"
var lock = sync.Mutex{}

func main() {
	listenAddr := os.Args[1]
	speakAddrs := save{}
	jsonUtil.LoadFromFile(serverFile, speakAddrs)

	var sPick string = ""
	for {
		lock.Lock()
		pickSpeak(&sPick, speakAddrs)
		go OnConnect(listenAddr, sPick)
	}

}

func pickSpeak(lastSpeak *string, sAddrs save) error {
	addrsCnt := len(sAddrs.Servers)
	if addrsCnt == 0 {
		return errors.New("No servers configured in " + serverFile)
	}

	if *lastSpeak == "" {
		this := sAddrs.Servers[0]
		*lastSpeak = this.IP + this.Port
		log.Println("Choosing " + *lastSpeak + " to serve")
		return nil
	}
	
	for i := 0; i < addrsCnt; i++ {
		this := sAddrs.Servers[i]
		if *lastSpeak == this.IP + this.Port{
			next := sAddrs.Servers[(i+1)%addrsCnt]
			*lastSpeak = next.IP + next.Port
			log.Println("Choosing " + *lastSpeak + " to serve")
			break
		}
	}

	return nil
}

func OnConnect(lAddr, sAddr string) error {

	lCon, sCon, er := func() (*net.Conn, *net.Conn, error) {
		defer lock.Unlock()
		log.Println("Open to connect on " + lAddr)

		l, er := connect.OpenConnection(lAddr)
		if er != nil {
			return nil, nil, er
		}
		log.Println("A client has connected to " + lAddr)

		s, er := connect.SeekConnection("", 5)
		if er != nil {
			c := *l; c.Close()
			return nil, nil, er
		}
		log.Println("The reverse proxy has connected to " + sAddr)

		return l, s, nil
	}()

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
