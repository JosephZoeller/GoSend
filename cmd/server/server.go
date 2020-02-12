package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/logutil"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

// Opens listening connections, then awaits a signal interruption to terminate.
func main() {
	var er error
	er = os.Mkdir("./downloads/", 0777)
	if er != nil && !os.IsExist(er) {
		log.Fatal("[Log Manager]: Failed to create downloads directory.")
	}

	log.Println("Server is checking for a connection with the Log Manager...")
	logConn, er = logutil.ConnectLog(logAddr)
	if er != nil {
		log.Println("Server did not connect with the Log Manager - ", er)
	} else {
		logutil.SendLog(logConn, fmt.Sprintf("Server connected with the Log Manager at [%s]", logAddr))
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

	logutil.SendLog(logConn, fmt.Sprintf("Opening Server connection at [%s]", transferAddress))
	conn, er := connect.OpenConnection(transferAddress)
	if er != nil {
		logutil.SendLog(logConn, "Connection failed - "+er.Error())
		return
	}
	c := *conn
	defer c.Close()
	logutil.SendLog(logConn, fmt.Sprintf("Server connection established at [%s]", transferAddress))

	EoFCnt := 0
	for {
		fHeader, er := transit.HeaderInbound(conn)
		if EoFCnt > 3 { // arbitrary
			logutil.SendLog(logConn, "End of session assumed. Closing connection to Server.")
			break
		} else if er == io.EOF {
			EoFCnt++
			continue
		} else if er != nil {
			logutil.SendLog(logConn, "Failed to receive file header - "+er.Error())
			continue
		}
		EoFCnt = 0

		c.SetReadDeadline(time.Now().Add(5000000000)) // 5 secconds

		oldname := fHeader.Filename
		fHeader.Filename = "./downloads/" + fHeader.Filename

		er = transit.FileInbound(fHeader, conn)
		if er != nil {
			logutil.SendLog(logConn, "Failed to receive file - "+er.Error())
		} else {
			logutil.SendLog(logConn, fmt.Sprintf("Success - Server received file %s. Server is awaiting next transmission...", oldname))
		}

		c.SetReadDeadline(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)) //unset
	}
}
