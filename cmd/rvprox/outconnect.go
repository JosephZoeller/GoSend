package main

import (
	"fmt"
	"net"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/logutil"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

// Creates a connection between the proxy and a server address.
func plugConnection(address string) (*net.Conn, error) {
	sCon, er := connect.SeekConnection(address, 5)
	if er != nil {
		return nil, er
	}
	logutil.SendLog(logConn, fmt.Sprintf("Proxy plugged into the Server at [%s]", address))

	return sCon, nil
}

// Closes all connections between the proxy and the server addresses.
func closeConnections() {
	logutil.SendLog(logConn, "Proxy is closing connections with the server.")
	for _, v := range outConns {
		c := *v
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
	logutil.SendLog(logConn, fmt.Sprintf("Proxy selected Server [%s] for data transmission.", c.RemoteAddr().String()))
	return lastConn
}
