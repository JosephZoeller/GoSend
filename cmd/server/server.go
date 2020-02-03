package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	if len(inAddrs) == 0 {
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
	for {
		log.Println("[Server serve]: Opening connection at " + transferAddress)
		conn, er := connect.OpenConnection(transferAddress)
		if er != nil {
			log.Println(logUtil.FormatError("Server serve", er))
			return
		}
		c := *conn
		defer c.Close()

		fHeader, er := transit.HeaderInbound(conn)
		if er != nil {
			log.Println(logUtil.FormatError("Server serve", er))
			return
		}

		er = transit.FileInbound(fHeader, conn)
		if er != nil {
			log.Println(logUtil.FormatError("Server serve", er))
			return
		}
	}
}
