package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type SaveFile struct {
	Thoughts []UserThought `json:"Thoughts"`
}

type UserThought struct {
	Date    string `json:"Date"`
	User    string `json:"User"`
	Thought string `json:"Thought"`
}

const savefilename string = "save.json"

func main() {
	hostSave()
}

func loadFile() (*SaveFile, error) {
	saves := SaveFile{}

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

func saveFile(s *SaveFile) error {
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

func hostSave() {
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
		up := UserThought{}

		log.Println("Upload Requested")
		if req.Method != "POST" {
			log.Println("Upload Rejected")
		} else {
			err := json.NewDecoder(req.Body).Decode(&up)

			if err != nil {
				log.Println(err)
			} else {
				AppendSave(up)
			}
		}
	})

	errorChan := make(chan error)
	fmt.Println("Listening on port 8080 (http)...")
	go func() {
		errorChan <- http.ListenAndServe(":8080", nil)
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
			log.Println("shutting down: ", sig)
			os.Exit(0)
		}
	}

}

func AppendSave(usTh UserThought) {

	saves, er := loadFile()
	if er != nil {
		log.Println(er)
	}

	saves.Thoughts = append(saves.Thoughts, usTh)

	er = saveFile(saves)
	if er != nil {
		log.Println(er)
	} else {
		log.Println("Content Updated")
	}

}
