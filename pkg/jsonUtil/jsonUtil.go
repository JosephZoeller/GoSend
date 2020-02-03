package jsonUtil

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/JosephZoeller/gmg/pkg/logUtil"
)

// SaveToFile attempts to marshal an interface into a json file.
// Accepts an interface to marshal and the desired filename.
func SaveToFile(fileName string, v interface{}) error {
	file, er := os.Create(fileName)
	if er != nil {
		return logUtil.FormatError("JsonUtil SaveToFile", er)
	}
	defer file.Close()

	b, er := json.MarshalIndent(v, "", "  ")
	if er != nil {
		return logUtil.FormatError("JsonUtil SaveToFile", er)
	}

	_, er = file.Write(b)
	if er != nil {
		return logUtil.FormatError("JsonUtil SaveToFile", er)
	}
	return nil
}

// Attempts to unmarshal a file into an interface.
// Accepts a filename and an interface for the unmarshalled content.
func LoadFromFile(fileName string, v interface{}) error {
	b, er := ioutil.ReadFile(fileName)
	if er != nil {
		return logUtil.FormatError("JsonUtil LoadFromFile", er)
	}
	er = json.Unmarshal(b, v)
	if er != nil {
		return logUtil.FormatError("JsonUtil LoadFromFile", er)
	}
	return nil
}

/*
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
*/

/*
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
*/