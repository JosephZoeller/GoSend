package transit

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"os"

	"github.com/JosephZoeller/gmg/pkg/logUtil"
)

// FileInbound Downloads the connection stream to a file.
// Requires a file header to determine the size and name of the file.
func FileInbound(fHead *fileHeader, con *net.Conn) error {
	c := *con

	fileCreate, er := os.Create(fHead.Filename)
	if er != nil {
		return logUtil.FormatError("Transit FileInbound", er)
	}

	for i := int64(0); i <= fHead.Blocks; i++ {
		_, er = io.CopyN(fileCreate, c, 1024)
		if er != nil {
			return logUtil.FormatError("Transit FileInbound", er)
		}
	}

	er = fileCreate.Truncate(fHead.Blocks*1024 + fHead.TailSize)
	if er != nil {
		return logUtil.FormatError("Transit FileInbound", er)
	}

	er = fileCreate.Close()
	if er != nil {
		return logUtil.FormatError("Transit FileInbound", er)
	}
	return nil
}

// HeaderInbound Downloads a file header from the connection stream.
func HeaderInbound(con *net.Conn) (*fileHeader, error) {
	c := *con
	tHeader := fileHeader{}

	jsonHeader := make([]byte, 1024)
	_, er := c.Read(jsonHeader)
	if er != nil {
		return &tHeader, logUtil.FormatError("Transit HeaderInbound", er)
	}

	jsonHeader = bytes.Trim(jsonHeader, "\x00")
	er = json.Unmarshal(jsonHeader, &tHeader)
	if er != nil {
		return &tHeader, logUtil.FormatError("Transit HeaderInbound", er)
	}

	return &tHeader, nil
}
