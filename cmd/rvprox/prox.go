package main

import (
	"errors"
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

var inAddrs []string
var outAddrs []string
var outConns = make([]*net.Conn, 0)
var lastConn *net.Conn

func init() {
	inAddrs, outAddrs = connect.ThroughArgs()
}

// Opens all proxy listeners, then awaits a signal interrupt to terminate.
func main() {

	if len(inAddrs) == 1 && inAddrs[0] == "" {
		log.Fatal(logUtil.FormatError("Proxy args", errors.New("No inbound addresses declared by user.")))
		return
	} else if len(outAddrs) == 1 && outAddrs[0] == "" {
		log.Fatal(logUtil.FormatError("Proxy args", errors.New("No outbound addresses declared by user.")))
		return
	}

	// How do you make sure that no one can directly access the server?
	// My solution is to plug the holes and keep them plugged until the end of time.
	plugAll()
	lastConn = outConns[len(outConns)-1]; defer closeConnections()

	log.Printf("Registered %d ports", len(inAddrs))
	for i, v := range inAddrs {
		go openListener(v)
		log.Printf("Proxy '%d' listening at '%s'", i+1, v)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan
}

// Gets all connections between the proxy and the server addresses.
func plugAll() {
	for _, v := range outAddrs {
		c, er := plugConnection(v)
		if er != nil {
			log.Fatal("Couldn't connect to an out address: " + v)
		}
		outConns = append(outConns, c)
	}
}

// Creates a connection between the proxy and a server address.
func plugConnection(address string) (*net.Conn, error) {
	sCon, er := connect.SeekConnection(address, 5)
	if er != nil {
		return nil, logUtil.FormatError("Proxy sendToAddress", er)
	}
	log.Println("[Proxy plugConnection]: The reverse proxy has plugged into " + address)

	return sCon, nil
}

// Closes all connections between the proxy and the server addresses.
func closeConnections() {
	for _, v := range outConns {
		c := *v
		c.Close()
	}
}

// Opens address to listen to and, upon connecting, select a speaking address and reroute the transmission.
func openListener(address string) {

	for {
		lCon, er := connect.OpenConnection(address)
		if er != nil {
			log.Println(logUtil.FormatError("Proxy OpenConnection", er))
			break
		}
		for {
			er = sendToAddress(lCon, pickSConn())
			if er == io.EOF {
				break
			}
			if er == nil {
				break
			}
			log.Println(logUtil.FormatError("Proxy openListener", er))
		}
		c := *lCon
		c.Close()
	}
}

// Reroutes the transmission being sent from the listening connection to the speaking connection.
// Attempts to speak to an address for 5 seconds before getting bored.
func sendToAddress(lCon, sCon *net.Conn) error {

	fHead, er := transit.PassHeader(lCon, sCon)
	if er != nil {
		return logUtil.FormatError("Proxy PassHeader", er)
	}

	er = transit.PassFile(fHead, lCon, sCon)
	if er != nil {
		log.Printf("[Proxy PassFile]: failed to pass %s (blocksize: %d, tailsize: %d)", fHead.Filename, fHead.Kilobytes, fHead.TailSize)
		return logUtil.FormatError("Proxy PassFile", er)
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

	c := *lastConn; log.Println("[Proxy pickSConn]: Connection chosen for transmission:", c.RemoteAddr())
	return lastConn
}
