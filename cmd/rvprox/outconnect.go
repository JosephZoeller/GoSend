package main

import (
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
	logUtil.SendLog(logConn, "Proxy plugged into the Server address at "+address)

	return sCon, nil
}

// Closes all connections between the proxy and the server addresses.
func closeConnections() {
	for _, v := range outConns {
		c := *v
		c.Close()
	}
}
