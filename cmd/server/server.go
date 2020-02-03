package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

var inAddrs []string

//var displayAddr string
func init() {
	inAddrs = connect.InArgs()
}

func main() {

	for i :=0; i < len(inAddrs); i++ {
		go serve(inAddrs[i])
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan
}

 func serve(transferAddr string) {
	 for {
		log.Println("Open connection at " + transferAddr)
		conn, er := connect.OpenConnection(transferAddr)
		if er != nil {
			log.Println("Get Session Error: ", er)
			return
		}
		c := *conn; defer c.Close()
	
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
	 }
 }