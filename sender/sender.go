package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// post save data for typetestd to process.
type UserThought struct {
	Date    string `json:"Date"`
	User    string `json:"User"`
	Thought string `json:"Thought"`
}

func main() {
	SendThought("Date Placeholder", "User Placeholder", "When we sleep, where do we go?")
	SendThought("Date Placeholder2", "User Placeholder2", "Why can't woodchucks chuck wood?")
}

func SendThought(d, u string, th string) {
	up := UserThought{
		Date:    d,
		User:    u,
		Thought: th,
	}

	js, err := json.Marshal(up)
	if err != nil {
		log.Println("[Save]: " + err.Error())
		return
	}

	res, err := http.Post("http://localhost:8080/Upload", "application/json", bytes.NewBuffer(js))
	if err != nil {
		log.Println("[Save]: " + err.Error())
	} else if res.StatusCode != 200 {
		log.Println("[Save]: " + res.Status)
	}
}
