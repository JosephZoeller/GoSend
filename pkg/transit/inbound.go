package transit

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"os"
)

func FileInbound(fHead *fileHeader, con *net.Conn) error {
	c := *con

	fileCreate, er := os.Create(fHead.Filename)
	if er != nil {
		return er
	}
	//buf := make([]byte, 1024)
	for i := int64(0); i <= fHead.Blocks; i++ {
		io.CopyN(fileCreate, c, 1024)
	}

	er = fileCreate.Truncate(fHead.Blocks*1024 + fHead.TailSize)
	if er != nil {
		return er
	}

	er = fileCreate.Close()
	if er != nil {
		return er
	}
	return nil
}

func HeaderInbound(con *net.Conn) (*fileHeader, error) {
	c := *con
	tHeader := fileHeader{}

	jsonHeader := make([]byte, 1024)
	_, er := c.Read(jsonHeader)
	if er != nil {
		return &tHeader, er
	}

	jsonHeader = bytes.Trim(jsonHeader, "\x00")
	er = json.Unmarshal(jsonHeader, &tHeader)
	if er != nil {
		return &tHeader, er
	}

	return &tHeader, nil
}
