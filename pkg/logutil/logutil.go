package logutil

import (
	"errors"
	"log"
	"net"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

// SendLog sends the log message to the log manager
func SendLog(sCon *net.Conn, msg string) {
	if sCon != nil {
		c := *sCon
		transit.LogOutbound(transit.MakeLogMsg(c.RemoteAddr().String(), msg), sCon)
	} /* else {
		log.Println(msg)
	}
	*/
	log.Println(msg)
}

// ConnectLog Dials the log manager
func ConnectLog(logAddr string) (*net.Conn, error) {
	if logAddr != "" {
		sCon, er := connect.SeekConnection(logAddr, 5)
		if er != nil {
			return nil, er
		}

		return sCon, nil
	}
	return nil, errors.New("N/A")
}
