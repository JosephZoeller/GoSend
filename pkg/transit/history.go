package transit

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const savefilename string = "save.json"

func saveToFile(s *saveFile) error {
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

func loadFromFile() (*saveFile, error) {
	saves := saveFile{}

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

func hostSave(addr string) {
	http.HandleFunc("/Display", func(res http.ResponseWriter, req *http.Request) {
		saves, er := loadFromFile()
		if er != nil {
			log.Println(er)
		}

		log.Println("Displaying Content")
		t, _ := template.ParseFiles("./web/tables.html")
		t.Execute(res, *saves)
	})

	errorChan := make(chan error)
	go func() {
		errorChan <- http.ListenAndServe(addr, nil)
		log.Printf("Listening on port :%s", addr)
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
			log.Printf("Shutting down: %s", sig)
			os.Exit(0)
		}
	}
}

func appendSave(tHeader *fileHeader) {

	saves, er := loadFromFile()
	if er != nil {
		log.Println(er)
	}

	saves.Files = append(saves.Files, *tHeader)

	er = saveToFile(saves)
	if er != nil {
		log.Println(er)
	} else {
		log.Println("Content Updated")
	}

}
