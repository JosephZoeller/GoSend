package main

import (
	"fmt"
	"io"
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
	log.Println("Proxy is checking for a connection with the Log Manager...")
	logConn, er = logUtil.ConnectLog(logAddr)
	if er != nil {
		log.Println("Proxy did not connect with the Log Manager - ", er)
	} else {
		logUtil.SendLog(logConn, fmt.Sprintf("Proxy connected with the Log Manager at [%s]", logAddr))
	}

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
				logUtil.SendLog(logConn, fmt.Sprintf("End of session assumed. Readying a new session at [%s]", c.RemoteAddr().String()))
				break
			} else if er == nil {
				logUtil.SendLog(logConn, fmt.Sprintf("File %s passing successful. Proxy will await the next transmission...", fName))
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

// Reroutes the transmission being sent from the listening connection to the speaking connection.
// Attempts to speak to an address for 5 seconds before getting bored.
func sendToAddress(lCon, sCon *net.Conn) (string, error) {

	fHead, er := transit.PassHeader(lCon, sCon)
	if er != nil {
		return "", er
	}
	fName := fHead.Filename
	er = transit.PassFile(fHead, lCon, sCon)
	if er != nil {
		return fName, er
	}

	return fName, nil
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
	logUtil.SendLog(logConn, fmt.Sprintf("Proxy selected Server [%s] for data transmission.", c.RemoteAddr().String()))
	return lastConn
}
