package main

import (
	"fmt"
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
	var er error
	var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	if fileSave {
		er = os.Mkdir("./logs/", 0777)
		if er != nil && !os.IsExist(er) {
			log.Fatal("[Log Manager]: Failed to create logs directory.")
		}
	}

	for i := 0; i < len(inAddrs); i++ {
		if fileSave {
			str := formatFilename(inAddrs[i])
			write, er := os.Create(str)
			if er == nil {
				logger = log.New(write, "", log.Ldate|log.Ltime)
				defer write.Close()
			} else {
				forkOutput(logger, fmt.Sprintf("[Log Manager]: Failed to create file for %s - %s", str, er.Error()))
			}
		}
		go logListen(inAddrs[i], logger)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan
}

func logListen(addrs string, logger *log.Logger) {

	forkOutput(logger, "[Log Manager]: Opening connection at address "+addrs)
	conn, er := connect.OpenConnection(addrs)
	if er != nil {
		forkOutput(logger, "[Log Manager]: Failed to open connection - "+er.Error())
		return
	}
	c := *conn
	defer c.Close()

	EoFCnt := 0
	for {
		if EoFCnt > 3 { // arbitrary
			forkOutput(logger, "[Log Manager]: Manager assumed End of Session, closing connection to Log Manager.")
			break
		}

		logmsg, er := transit.LogInbound(conn)
		if er == io.EOF {
			EoFCnt++
			continue
		} else if er != nil {
			forkOutput(logger, "[Log Manager]: Failed to receive log message - "+er.Error())
			continue
		}
		EoFCnt = 0

		forkOutput(logger, logmsg.String())
	}
}

func forkOutput(logger *log.Logger, msg string) {
	logger.Println(msg)
	if logger.Writer() != os.Stdout {
		log.Println(msg)
	}
}

func formatFilename(address string) string {
	return "./logs/" + strings.ReplaceAll(address, ":", "_") + ".log"
}
