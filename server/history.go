package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	structs "github.com/JosephZoeller/gmg/util"
)

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

func appendSave(tHeader *structs.FileHeader) {

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
