package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"time"

	structs "github.com/JosephZoeller/gmg/Util"
)

var user string = os.Getenv("USER")

func main() {
	conn, er := EstablishConnection(os.Args[1], 5)
	defer conn.Close()
	if er != nil {
		log.Println(er)
		return
	}

	thought, er := ReadThought()
	if er != nil {
		log.Println(er)
	}
	er = SendThought(MakeBody(thought), conn)
	if er != nil {
		log.Println(er)
	}
	/*
		for {
			thought, er := ReadThought()
			if er != nil {
				log.Println(er)
			} else {
				SendThought(MakeBody(thought), conn)
			}
		}
	*/

}

func EstablishConnection(port string, timeout int) (net.Conn, error) {
	for i := 0; i <= timeout; i++ {
		c, er := net.Dial("tcp", "localhost:"+port)
		if er == nil {
			return c, nil
		}
		time.Sleep(time.Second)
	}
	return nil, errors.New("[Establish Connection]: Connection Timed Out")
}

func MakeBody(thought string) structs.ThoughtBody {
	return structs.ThoughtBody{
		User:    user,
		Date:    time.Now().Format("Jan/2/2006"),
		Thought: TrimThought(thought),
	}
}

func ReadThought() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	thought, er := reader.ReadString('\n')
	if er != nil {
		return "", er
	} else {
		return thought, nil
	}
}

func TrimThought(thought string) string {
	return strings.TrimSuffix(thought, "\n")
}

func SendThought(tBody structs.ThoughtBody, conn net.Conn) error { // send packaged thought
	jsonBody, err := json.Marshal(tBody)
	if err != nil {
		return err
	}

	tHeader := structs.ThoughtHeader{
		Size: (len(jsonBody))/1024 + 1,
	}
	jsonHeader, err := json.Marshal(tHeader)
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)
	copy(buf, jsonHeader)
	_, err = conn.Write(buf)

	if err != nil {
		return err
	}

	for i := 1; i < tHeader.Size; i++ {
		buf := jsonBody[i-1 : i*1024]
		_, err = conn.Write(buf)
		if err != nil {
			return err
		}
	}

	buf = make([]byte, 1024)
	copy(buf, jsonBody)
	_, err = conn.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
