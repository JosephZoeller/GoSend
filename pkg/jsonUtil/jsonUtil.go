package jsonUtil

import (
	"encoding/json"
	"os"
)

func SaveToFile(fileName string, v interface{}) error {
	file, er := os.Create(fileName)
	if er != nil {
		return er
	}
	defer file.Close()

	b, er := json.MarshalIndent(v, "", "  ")
	if er != nil {
		return er
	}

	_, er = file.Write(b)
	return er
}

func LoadFromFile(fileName string, v interface{}) error {
	file, er := os.Open(fileName)
	if er != nil {
		return er
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(v)
}
