package transit

import (
	"os"
	"runtime"
	"time"
)

type saveFile struct {
	Files []fileHeader `json:"Files"`
}

type fileHeader struct {
	Filename  string `json:"Filename"`
	User      string `json:"User"`
	Date      string `json:"Date"`
	AuthToken string `json:"Authentication"`
	Kilobytes int64  `json:"Blocks"`
	TailSize  int64  `json:"Tail"`
}

// creates a header which breifs the transmission reciever on what to expect from the incoming file.
func MakeHeader(file *os.File) (*fileHeader, error) {
	fstat, er := file.Stat()
	if er != nil {
		return &fileHeader{}, er
	}

	return &fileHeader{
		User:      getDefaultName(),
		Date:      time.Now().Format("Jan/2/2006"),
		Kilobytes: fstat.Size() / 1024,
		TailSize:  fstat.Size() % 1024,
		Filename:  fstat.Name(),
	}, nil
}

// Gets the environment username. Operating System-dependant.
func getDefaultName() string {
	var userEnvVar string
	if runtime.GOOS == "windows" {
		userEnvVar = os.Getenv("USERNAME")
	} else if runtime.GOOS == "linux" {
		userEnvVar = os.Getenv("USER")
	} else {
		userEnvVar = "User"
	}
	return userEnvVar
}
