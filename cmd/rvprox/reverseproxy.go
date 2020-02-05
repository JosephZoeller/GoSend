package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/logUtil"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

// Opens all proxy listeners, then awaits a signal interrupt to terminate.
func main() {
	var er error
	logConn, er = logUtil.ConnectLog(logAddr)
	if er != nil {
		log.Println("Proxy failed to connect with Log Manager - ", er)
	} else {
		logUtil.SendLog(logConn, " Proxy connected with Log Manager")
	}

	connectServers(); defer closeConnections()
	logUtil.SendLog(logConn, fmt.Sprintf("Proxy registered all %d server addresses", len(inAddrs)))
	lastConn = outConns[len(outConns)-1]

	for _, v := range inAddrs {
		go openListener(v)
		logUtil.SendLog(logConn, fmt.Sprintf("Proxy listening at '%s'", v))
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
			logUtil.SendLog(logConn, fmt.Sprintf("Proxy failed to connect with Server at address %s - %s", v, er.Error()))
			os.Exit(1)
		}
		outConns = append(outConns, c)
	}
}

// Opens address to listen to and, upon connecting, select a speaking address and reroute the transmission.
func openListener(address string) {

	for {
		lCon, er := connect.OpenConnection(address)
		if er != nil {
			logUtil.SendLog(logConn, fmt.Sprintf("Proxy failed to connect with Client on address %s - %s", address, er.Error()))
			break
		}

		c := *lCon; c.Close()
		for {
			er = sendToAddress(lCon, pickSConn())
			if er == nil {
				break
			}
			s := *lastConn
			logUtil.SendLog(logConn, fmt.Sprintf("Proxy failed to pass data from Client %s to Server %s - %s", c.RemoteAddr().String(), s.RemoteAddr().String(), er.Error()))
		}
	}
}

// Reroutes the transmission being sent from the listening connection to the speaking connection.
// Attempts to speak to an address for 5 seconds before getting bored.
func sendToAddress(lCon, sCon *net.Conn) error {

	fHead, er := transit.PassHeader(lCon, sCon)
	if er != nil {
		return er
	}

	er = transit.PassFile(fHead, lCon, sCon)
	if er != nil {
		return er
	}

	return nil
}

// Load balancer (if you can call it that). Selects (via round-robin) an address to speak to from the CLI argument.
func pickSConn() *net.Conn {
	connCnt := len(outConns)

	for i := 0; i < connCnt; i++ {
		if lastConn == outConns[i] {
			lastConn = outConns[(i+1)%connCnt]
			break
		}
	}

	c := *lastConn
	logUtil.SendLog(logConn, fmt.Sprintf("Proxy selected Server %s for data transmission.", c.LocalAddr().String()))
	return lastConn
}
