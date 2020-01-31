package main

import (
	"bytes"
	"encoding/json"
	//"fmt"
	//"bufio"
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
	ln, er := net.Listen("tcp", ":"+port)
	if er != nil {
		log.Println(er)
		return
	}
	conn, er := ln.Accept()
	if er != nil {
		log.Println(er)
		return
	}
	defer conn.Close()

	tHeader := structs.ThoughtHeader{}
	jsonHeader := make([]byte, 1024)
	_, er = conn.Read(jsonHeader)
	jsonHeader = bytes.Trim(jsonHeader, "\x00")
	if er != nil {
		log.Println("Header Read Error: ", er)
		return
	}

	er = json.Unmarshal(jsonHeader, &tHeader)
	if er != nil {
		log.Println("Unmarshal Error: ", er, "jsonHeader: ", string(jsonHeader))
		return
	}

	jsonBody := make([]byte, 0)
	for i := 0; i < tHeader.Size; i++ {
		buf := make([]byte, 1024)
		n, er := conn.Read(buf)
		if er != nil {
			log.Println("Read Error: ", er)
			return
		}
		jsonBody = append(jsonBody, buf[:n]...)
	}
	jsonBody = bytes.Trim(jsonBody, "\x00")

	/*
		bufio.NewWriter(os.Stdout).WriteString(string(jsonBody))
		bufio.NewReader(os.Stdin).ReadLine()
	*/

	tBody := structs.ThoughtBody{}
	log.Println(string(jsonBody))
	er = json.Unmarshal(jsonBody, &tBody)
	if er != nil {
		log.Println("Unmarshal Error: ", er, "jsonBody: ", string(jsonBody))
		return
	}
	AppendSave(tBody)
}

/*
var ConnSig chan string = make(chan string)

func main() {
	ln, _ := net.Listen("tcp", ":8080")
	for { // will create a whole bunch of sessions unless a blocker is introduced
		fmt.Println("Pre-Session...")
		go Session(ln)
		fmt.Println(<-ConnSig)
	}
}

func Session(ln net.Listener) {
	connection, _ := ln.Accept()
	ConnSig <- "connection est."
	defer connection.Close()
	for {
		buf := make([]byte, 1024) // 1024 is the standard
		connection.Read(buf)
		fmt.Println(string(buf))
	}
}
*/

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

	http.HandleFunc("/Upload", func(res http.ResponseWriter, req *http.Request) {
		tBody := structs.ThoughtBody{}

		log.Println("Upload Requested")
		err := json.NewDecoder(req.Body).Decode(&tBody)

		if err != nil {
			log.Println(err)
		} else {
			AppendSave(tBody)
		}
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

func AppendSave(tBody structs.ThoughtBody) {

	saves, er := loadFile()
	if er != nil {
		log.Println(er)
	}

	saves.Thoughts = append(saves.Thoughts, tBody)

	er = saveFile(saves)
	if er != nil {
		log.Println(er)
	} else {
		log.Println("Content Updated")
	}

}
