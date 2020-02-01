package main

import (
	"net"
	"os"
	"time"
	"encoding/json"
	structs "github.com/JosephZoeller/gmg/util"
)

// Sends a pre-file header (max 1kb) to the tcp connection
func headerOutbound(tHeader *structs.FileHeader, conn *net.Conn) error {
	c := *conn

	jsonHeader, er := json.Marshal(*tHeader)
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
func fileOutbound(fileIn *os.File, conn *net.Conn, tHeader *structs.FileHeader) error {
	c := *conn
	buf := make([]byte, 1024)
	for i := int64(1); i <= tHeader.Blocks; i++ {
		if tHeader.Blocks == i {		// final bytes, only TailSize and empty buffer bytes remaining
			buf = make([]byte, 1024) 	// clears out old bytes in the buffer
		}

		er := func() error {
			_, er := fileIn.Read(buf)
			if er != nil {
				return er
			}

			_, er = c.Write(buf)
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

// creates a pre-file header to communicate with the server.
func makeHeader(fileIn *os.File) (*structs.FileHeader, error) {
	fstat, er := fileIn.Stat()
	if er != nil {
		return &structs.FileHeader{}, er
	}

	return &structs.FileHeader{
		User:     user,
		Date:     time.Now().Format("Jan/2/2006"),
		Blocks:   fstat.Size()/1024 + 1,
		TailSize: int(fstat.Size() % 1024),
		Filename: fstat.Name(),
	}, nil
}
