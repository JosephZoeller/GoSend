package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	structs "github.com/JosephZoeller/gmg/Util"
)

const savefilename string = "save.json"

var server string
var port string

func main() {

	server = os.Args[1]
	port = os.Args[2]
	go hostSave("8081")
	conn, er := GetSession(port)
	if er != nil {
		log.Println("Get Session Error: ", er)
		return
	}
	c := *conn; defer c.Close()

	tHeader, er := HeaderInbound(conn)
	if er != nil {
		log.Println("Create Error: ", er)
		return
	}

	fileCreate, er := os.Create(tHeader.Filename)
	if er != nil {
		log.Println("Create Error: ", er)
		return
	}
	defer fileCreate.Close()

	if FileInbound(tHeader, conn, fileCreate) != nil {
		return
	}
	AppendSave(tHeader)
}

func FileInbound(tHeader *structs.FileHeader, con *net.Conn, fileOut *os.File) error {
	c := *con
	buf := make([]byte, 1024)

	for i := int64(1); i <= tHeader.Blocks; i++ {
		if tHeader.Blocks == i {
			buf = make([]byte, 1024)
		}

		er := func() error {
			_, er := c.Read(buf)
			if er != nil {
				return er
			}

			if tHeader.Blocks == i {
				_, er = fileOut.Write(buf[:tHeader.TailSize])
			} else {
				_, er = fileOut.Write(buf)
			}
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

func HeaderInbound(con *net.Conn) (*structs.FileHeader, error) {
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

func GetSession(port string) (*net.Conn, error){
	ln, er := net.Listen("tcp", ":"+port)
	if er != nil {
		log.Println(er)
		return nil, er
	}

	conn, er := ln.Accept()
	if er != nil {
		log.Println(er)
		return nil, er
	}

	return &conn, er
}

func loadFile() (*structs.SaveFile, error) {
	saves := structs.SaveFile{}

	file, er := os.Open(savefilename)
	if er != nil {
		return &saves, er
	} else {
		defer file.Close()
	}

	er = json.NewDecoder(file).Decode(&saves)
	if er != nil {
		return &saves, er
	}

	return &saves, nil
}

func saveFile(s *structs.SaveFile) error {
	file, er := os.Create(savefilename)
	if er != nil {
		return er
	}
	defer file.Close()

	en := json.NewEncoder(file)
	en.SetIndent("", "  ")
	er = en.Encode(*s)
	if er != nil {
		return er
	}
	return nil
}

func hostSave(portNum string) {
	http.HandleFunc("/Display", func(res http.ResponseWriter, req *http.Request) {
		saves, er := loadFile()
		if er != nil {
			log.Println(er)
		}

		log.Println("Displaying Content")
		t, _ := template.ParseFiles("./web/tables.html")
		t.Execute(res, *saves)
	})

	errorChan := make(chan error)
	go func() {
		errorChan <- http.ListenAndServe(":"+portNum, nil)
		log.Printf("%s listening on port :%s", server, portNum)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	for {
		select {
		case err := <-errorChan:
			if err != nil {
				log.Fatalln(err)
			}

		case sig := <-signalChan:
			log.Printf("Server %s shutting down: %s", server, sig)
			os.Exit(0)
		}
	}

}

func AppendSave(tHeader *structs.FileHeader) {

	saves, er := loadFile()
	if er != nil {
		log.Println(er)
	}

	saves.Files = append(saves.Files, *tHeader)

	er = saveFile(saves)
	if er != nil {
		log.Println(er)
	} else {
		log.Println("Content Updated")
	}

}
