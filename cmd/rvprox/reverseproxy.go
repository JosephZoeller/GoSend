package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/logUtil"
)

// Opens all proxy listeners, then awaits a signal interrupt to terminate.
func main() {
	var er error
	log.Println("Proxy is checking for a connection with the Log Manager...")
	logConn, er = logUtil.ConnectLog(logAddr)
	if er != nil {
		log.Println("Proxy did not connect with the Log Manager - ", er)
	} else {
		logUtil.SendLog(logConn, fmt.Sprintf("Proxy connected with the Log Manager at [%s]", logAddr))
	}

	logUtil.SendLog(logConn, fmt.Sprintf("Proxy is attempting to connect with all server addresses..."))
	connectServers()
	defer closeConnections()
	logUtil.SendLog(logConn, fmt.Sprintf("Proxy registered all %d server addresses", len(outAddrs)))
	lastConn = outConns[len(outConns)-1]

	for _, v := range inAddrs {
		go openListener(v)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan
}

// Gets all connections between the proxy and the server addresses.
func connectServers() { // TODO have servers connect to the proxy rather than have the proxy connect to the servers.
	for _, v := range outAddrs {
		c, er := plugConnection(v)
		if er != nil {
			logUtil.SendLog(logConn, fmt.Sprintf("Proxy failed to connect with a Server at address [%s] - %s", v, er.Error()))
			os.Exit(1)
		}
		outConns = append(outConns, c)
	}
}

// Opens address to listen to and, upon connecting, select a speaking address and reroute the transmission.
func openListener(address string) {
	logUtil.SendLog(logConn, fmt.Sprintf("Proxy listening at [%s]", address))
	for {
		lCon, er := connect.OpenConnection(address)
		if er != nil {
			logUtil.SendLog(logConn, fmt.Sprintf("Proxy failed to connect with Client on address [%s] - %s", address, er.Error()))
			break
		}
		c := *lCon
		logUtil.SendLog(logConn, fmt.Sprintf("Proxy connection established at [%s]", address))

		EoFCnt := 0
		fName, er := sendToAddress(lCon, pickSConn())
		for { // loops for each attempt at passing a file
			if EoFCnt > 3 {
				logUtil.SendLog(logConn, fmt.Sprintf("Proxy assumed End of Session. Readying a new session at [%s]", c.RemoteAddr().String()))
				break
			} else if er == nil {
				logUtil.SendLog(logConn, fmt.Sprintf("Success - Proxy passed file %s. Proxy will await the next transmission...", fName))
				break
			} else if er == io.EOF {
				EoFCnt++
				_, er = sendToAddress(lCon, lastConn)
				continue
			}

			s := *lastConn
			logUtil.SendLog(logConn, fmt.Sprintf("Proxy failed to pass file '%s' from Client [%s] to Server [%s] - %s", fName, c.RemoteAddr().String(), s.RemoteAddr().String(), er.Error()))
			break
		}
		c.Close()
	}
}
