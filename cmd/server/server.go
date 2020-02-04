package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/logUtil"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

var inAddrs []string

func init() {
	inAddrs = connect.InArgs()
}

// Opens listening connections, then awaits a signal interruption to terminate.
func main() {
	if len(inAddrs) == 1 && inAddrs[0] == "" {
		log.Println(logUtil.FormatError("Server args", errors.New("No inbound addresses declared by user.")))
		return
	}

	for i := 0; i < len(inAddrs); i++ {
		go serve(inAddrs[i])
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan
}

// Opens a connection on the transferAddress and, upon connecting, receives the transmission data.
func serve(transferAddress string) {

	log.Println("[Server serve]: Opening connection at " + transferAddress)
	conn, er := connect.OpenConnection(transferAddress)
	if er != nil {
		log.Println(logUtil.FormatError("Server OpenConnection", er))
		return
	}
	c := *conn
	defer c.Close()

	for {
		fHeader, er := transit.HeaderInbound(conn)
		if er != nil {
			log.Println(logUtil.FormatError("Server HeaderInbound", er))
			continue
		}

		c.SetReadDeadline(time.Now().Add(5000000000)) // 5 secconds
		er = transit.FileInbound(fHeader, conn)
		if er != nil {
			log.Println(logUtil.FormatError("Server FileInbound", er))
		}
		c.SetReadDeadline(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)) //unset
	}
}
