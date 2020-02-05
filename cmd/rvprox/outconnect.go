package main

import (
	"fmt"
	"net"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/logUtil"
)

// Creates a connection between the proxy and a server address.
func plugConnection(address string) (*net.Conn, error) {
	sCon, er := connect.SeekConnection(address, 5)
	if er != nil {
		return nil, er
	}
	logUtil.SendLog(logConn, fmt.Sprintf("Proxy plugged into the Server at [%s]", address))

	return sCon, nil
}

// Closes all connections between the proxy and the server addresses.
func closeConnections() {
	logUtil.SendLog(logConn, "Proxy is closing connections with the server.")
	for _, v := range outConns {
		c := *v
		c.Close()
	}
}
