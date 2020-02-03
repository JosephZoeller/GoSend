package transit

import (
	"encoding/json"
	"net"
	"os"
)

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
	return nil
}

// FileOutbound uploads a file to the connection stream.
func FileOutbound(fileOut *os.File, conn *net.Conn) error {
	c := *conn
	fstat, er := fileOut.Stat()
	if er != nil {
		return er
	}

	buf := make([]byte, 1024)
	for i := int64(0); i <= fstat.Size()/1024; i++ {
		_, er = fileOut.Read(buf)
		c.Write(buf)
	}

	return nil
}
