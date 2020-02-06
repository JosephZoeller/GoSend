package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

// hook the proxy and servers to the logging manager. Connections are infinite in order to leverage security with project deadlines.
func main() {
	var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	for i := 0; i < len(inAddrs); i++ {
		if fileSave {
			str := formatFilename(inAddrs[i])
			write, er := os.Create(str)
			if er == nil {
				logger = log.New(write, "", log.Ldate|log.Ltime)
				defer write.Close()
			} else {
				logger.Printf("[Log Manager]: Failed to create file for %s - %s", str, er.Error())
			}
		}
		go logListen(inAddrs[i], logger)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan
}

func logListen(addrs string, logger *log.Logger) {

	logger.Println("[Log Manager]: Opening connection at address " + addrs)
	conn, er := connect.OpenConnection(addrs)
	if er != nil {
		logger.Println("[Log Manager]: Failed to open connection - " + er.Error())
		return
	}
	c := *conn
	defer c.Close()

	EoFCnt := 0
	for {
		if EoFCnt > 3 { // arbitrary
			logger.Println("[Log Manager]: Manager assumed End of Session, closing connection to Log Manager.")
			break
		}

		logmsg, er := transit.LogInbound(conn)
		if er == io.EOF {
			EoFCnt++
			continue
		} else if er != nil {
			logger.Println("[Log Manager]: Failed to receive log message - " + er.Error())
			continue
		}
		EoFCnt = 0

		logger.Println(logmsg.String())
	}
}

func displayPlusSave(logger *log.Logger, msg string) {
	logger.Println(msg)
	if logger.Writer() != os.Stdin {
		log.Println(msg)
	}
}

func formatFilename(address string) string {
	return strings.ReplaceAll(address, ":", "_") + ".log"
}
