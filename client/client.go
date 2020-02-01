package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"time"

	structs "github.com/JosephZoeller/gmg/Util"
)

var user string = os.Getenv("USER")

func main() {

	fileIn, er := os.Open(os.Args[1])
	if er != nil {
		log.Println(er)
		return
	}

	conn, er := EstablishConnection(os.Args[2], 5)
	if er != nil {
		log.Println(er)
		return
	}
	c := *conn
	defer c.Close()

	tHeader, er := MakeHeader(fileIn)
	if er != nil {
		log.Println(er)
		return
	}
	er = HeaderOutbound(tHeader, conn)
	if er != nil {
		log.Println(er)
		return
	}

	er = FileOutbound(fileIn, conn, tHeader)
	if er != nil {
		log.Println(er)
	}
}

func EstablishConnection(port string, timeout int) (*net.Conn, error) {
	for i := 0; i <= timeout; i++ {
		c, er := net.Dial("tcp", "localhost:"+port)
		if er == nil {
			log.Println("[Establish Connection]: Connected")
			return &c, nil
		}
		time.Sleep(time.Second)
	}
	return nil, errors.New("[Establish Connection]: Connection Timed Out")
}

func MakeHeader(fileIn *os.File) (*structs.FileHeader, error) {
	fstat, er := fileIn.Stat()
	if er != nil {
		return &structs.FileHeader{}, er
	}

	return &structs.FileHeader{
		User:     user,
		Date:     time.Now().Format("Jan/2/2006"),
		Blocks:   fstat.Size()/1024 + 1,
		TailSize: fstat.Size() % 1024,
		Filename: fstat.Name(),
	}, nil
}

func HeaderOutbound(tHeader *structs.FileHeader, conn *net.Conn) error {
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

func FileOutbound(fileIn *os.File, conn *net.Conn, tHeader *structs.FileHeader) error {
	c := *conn
	buf := make([]byte, 1024)
	for i := int64(1); i <= tHeader.Blocks; i++ {
		if tHeader.Blocks == i {
			buf = make([]byte, 1024) // clears out old bytes in the buffer
		}

		er := func() error {
			_, er := fileIn.Read(buf)
			//log.Println(buf)
			if er != nil {
				//return err
			}

			_, er = c.Write(buf)
			if er != nil {
				//return err
			}

			return nil
		}()

		if er != nil {
			return er
		}
	}
	return nil
}
