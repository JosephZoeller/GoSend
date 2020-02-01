package main

import (
	"bytes"
	"encoding/json"
	"net"
	"os"

	structs "github.com/JosephZoeller/gmg/util"
)

// Anticipates a file. Reads (1024 * blocks) bytes, and writes (1024 * (blocks-1) + len(file tail)) bytes to a file.
func fileInbound(tHeader *structs.FileHeader, con *net.Conn, fileOut *os.File) error {
	c := *con
	buf := make([]byte, 1024)
	writelen := 1024

	for i := int64(1); i <= tHeader.Blocks; i++ {
		if tHeader.Blocks == i { // final bytes, only TailSize remaining.
			writelen = tHeader.TailSize
		}

		er := func() error {
			_, er := c.Read(buf)
			if er != nil {
				return er
			}

			_, er = fileOut.Write(buf[:writelen])
			if er != nil {
				return er
			}
			return nil
		}()

		if er != nil {
			return er
		}
	}
	return nil
}

// Anticipates a pre-file header in json format. The header (max 1kb) contains the filename, size (1024*(blocks-1) + tail size) sender and date.
func headerInbound(con *net.Conn) (*structs.FileHeader, error) {
	c := *con
	tHeader := structs.FileHeader{}

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
