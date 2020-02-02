package transit

import (
	"encoding/json"
	"net"
	"os"
)

// Sends a pre-file header (max 1kb) to the tcp connection
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

// Reads in bytes from a file as it streams them to the given connection.
// Reads len(file) bytes from the file and writes len(file) + (1024 - len(file)%1024) bytes to the connection.
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
