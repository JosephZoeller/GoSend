package transit

import (
	"encoding/json"
	"io"

	//"log"
	"net"
	"os"
)

// FileOutbound uploads a file to the connection stream.
func FileOutbound(fileOut *os.File, conn *net.Conn) error {
	c := *conn
	fstat, er := fileOut.Stat()
	if er != nil {
		return er
	}

	kb := fstat.Size() / 1024
	tail := fstat.Size() % 1024
	sendSize := int64(1024)

	for i := int64(0); i <= kb; i++ {
		if i == kb {
			sendSize = tail
		}

		_, er := io.CopyN(c, fileOut, sendSize)
		//log.Println(n)
		if er != nil {
			return er
		}

	}
	//log.Printf("[Outbound File]: Successfully sent %s file", fstat.Name())
	return nil
}

// HeaderOutbound uploads a file header to the connection stream.
func HeaderOutbound(fHead *fileHeader, conn *net.Conn) error {
	c := *conn

	jsonHeader, er := json.Marshal(fHead)
	if er != nil {
		return er
	}

	buf := make([]byte, 1024)
	copy(buf, jsonHeader)
	_, er = c.Write(buf)
	if er != nil {
		return er
	}

	//log.Printf("[Outbound Header]: Successfully sent %s header", fHead.Filename)
	return nil
}

// LogOutbound uploads a log to the connection stream.
func LogOutbound(msg *logMsg, conn *net.Conn) error {
	c := *conn

	jsonHeader, er := json.Marshal(msg)
	if er != nil {
		return er
	}

	buf := make([]byte, 1024)
	copy(buf, jsonHeader)
	_, er = c.Write(buf)
	if er != nil {
		return er
	}

	//log.Printf("[Outbound Log]: Successfully sent %s log message.", msg.SenderID)
	return nil
}
