package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/logUtil"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

// Opens listening connections, then awaits a signal interruption to terminate.
func main() {
	var er error
	logConn, er = logUtil.ConnectLog(logAddr)
	if er != nil {
		log.Println("[Proxy Log Connect]:", er)
	} else {
		logUtil.SendLog(logConn, " The reverse proxy has plugged into "+logAddr)
	}
	for i := 0; i < len(inAddrs); i++ {
		go serverListen(inAddrs[i])
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan
}

// Opens a connection on the transferAddress and, upon connecting, receives the transmission data.
func serverListen(transferAddress string) {

	log.Println("[Server Listen]: Opening connection to server at: " + transferAddress)
	conn, er := connect.OpenConnection(transferAddress)
	if er != nil {
		logUtil.SendLog(logConn, "Connection failed - "+er.Error())
		return
	}
	c := *conn
	defer c.Close()

	for {
		fHeader, er := transit.HeaderInbound(conn)
		if er != nil {
			logUtil.SendLog(logConn, "Failed to recieve file header - "+er.Error())
			continue
		}

		c.SetReadDeadline(time.Now().Add(5000000000)) // 5 secconds
		er = transit.FileInbound(fHeader, conn)
		if er != nil {
			logUtil.SendLog(logConn, "Failed to receive file - "+er.Error())
		}
		c.SetReadDeadline(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)) //unset
	}
}
