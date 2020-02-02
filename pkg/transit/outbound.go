package transit

import (
	"encoding/json"
	"net"
	"os"
	"runtime"
	"time"
)

// Sends a pre-file header (max 1kb) to the tcp connection
func HeaderOutbound(fileOut *os.File, conn *net.Conn) error {
	c := *conn
	fHead, er := makeHeader(fileOut)
	if er != nil {
		return er
	}

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
	fHead, er := makeHeader(fileOut)
	if er != nil {
		return er
	}

	buf := make([]byte, 1024)
	for i := int64(0); i <= fHead.Blocks; i++ {
		_, er = fileOut.Read(buf)
		c.Write(buf)
	}

	return nil
}

// creates a pre-file header to communicate with the server.
func makeHeader(fileOut *os.File) (*fileHeader, error) {
	fstat, er := fileOut.Stat()
	if er != nil {
		return &fileHeader{}, er
	}

	return &fileHeader{
		User:     getDefaultName(),
		Date:     time.Now().Format("Jan/2/2006"),
		Blocks:   fstat.Size() / 1024,
		TailSize: fstat.Size() % 1024,
		Filename: fstat.Name(),
	}, nil
}

func getDefaultName() string {
	var userEnvVar string
	if runtime.GOOS == "windows" {
		userEnvVar = os.Getenv("USERNAME")
	} else if runtime.GOOS == "linux" {
		userEnvVar = os.Getenv("USER")
	} else {
		userEnvVar = "User"
	}
	return userEnvVar
}
