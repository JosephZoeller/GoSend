package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

//  hook the proxy and servers to the logging manager. Connections are infinite in order to leverage security with project deadlines.
func main() {
	for i := 0; i < len(inAddrs); i++ {
		go logListen(inAddrs[i])
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan
}

func logListen(addrs string) {

	log.Println("[Log Manager]: Opening connection at address " + addrs)
	conn, er := connect.OpenConnection(addrs)
	if er != nil {
		log.Println("[Log Manager]: Failed to open connection - " + er.Error())
		return
	}
	c := *conn
	defer c.Close()

	EoFCnt := 0
	for {
		if EoFCnt > 3 { // arbitrary
			log.Println("[Log Manager]: End of File, closing connection to Log Manager.")
			break
		}

		logmsg, er := transit.LogInbound(conn)
		if er == io.EOF {
			EoFCnt++
			continue
		} else if er != nil {
			log.Println("[Log Manager]: Failed to recieve log message - " + er.Error())
			continue
		}
		EoFCnt = 0

		log.Println(logmsg.String())
	}
}
