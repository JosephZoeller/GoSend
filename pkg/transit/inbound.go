package transit

import (
	"bytes"
	"encoding/json"
	"io"

	//"path/filepath"

	//"log"
	"net"
	"os"
)

// FileInbound Downloads the connection stream to a file.
// Requires a file header to determine the size and name of the file.
func FileInbound(fHead *fileHeader, con *net.Conn) error {
	c := *con

	fileCreate, er := os.Create(fHead.Filename)
	if er != nil {
		return er
	}
	kb := fHead.Kilobytes
	tail := fHead.TailSize
	sendSize := int64(1024)

	for i := int64(0); i <= kb; i++ {
		if i == kb {
			sendSize = tail
		}
		_, er := io.CopyN(fileCreate, c, sendSize)
		//log.Println(n)
		if er != nil {
			return er
		}
	}

	//log.Printf("[Inbound File]: Successfully received %s file. Trimming tail...", fHead.Filename)
	er = fileCreate.Truncate(fHead.Kilobytes*1024 + fHead.TailSize) //TODO check if tail even needs trimming. is this deprecated?
	if er != nil {
		return er
	}

	er = fileCreate.Close()
	if er != nil {
		return er
	}

	//log.Printf("[Inbound File]: Successfully written file %s.", fHead.Filename)
	return nil
}

// HeaderInbound Downloads a file header from the connection stream.
func HeaderInbound(con *net.Conn) (*fileHeader, error) {
	c := *con
	fHead := fileHeader{}

	jsonHeader := make([]byte, 1024)
	_, er := c.Read(jsonHeader)
	if er != nil {
		return &fHead, er
	}

	jsonHeader = bytes.Trim(jsonHeader, "\x00")
	er = json.Unmarshal(jsonHeader, &fHead)
	if er != nil {
		return &fHead, er
	}

	//log.Printf("[Inbound Header]: Retrieved header for %s.", fHead.Filename)
	return &fHead, nil
}

// LogInbound acquires a log message from a connection
func LogInbound(con *net.Conn) (*logMsg, error) {
	c := *con
	msg := logMsg{}

	jsonHeader := make([]byte, 1024)
	_, er := c.Read(jsonHeader)
	if er != nil {
		return &msg, er
	}

	jsonHeader = bytes.Trim(jsonHeader, "\x00")
	er = json.Unmarshal(jsonHeader, &msg)
	if er != nil {
		return &msg, er
	}

	//log.Printf("[Inbound Log]: Retrieved message from %s.", msg.SenderID)
	return &msg, nil
}
