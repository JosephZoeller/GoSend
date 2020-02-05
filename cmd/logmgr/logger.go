package main

import (
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

	log.Println("[Log Listen]: Opening connection to log at address " + addrs)
	conn, er := connect.OpenConnection(addrs)
	if er != nil {
		log.Println("[Log Connection]:", er)
		return
	}
	c := *conn
	defer c.Close()

	for {
		logmsg, er := transit.LogInbound(conn)
		if er != nil {
			log.Println("[Log Inbound]:", er)
			continue
		}
		log.Println(logmsg.String())
	}
}
