package transit

import (
	"fmt"
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
	Kilobytes int64  `json:"KB"`
	TailSize  int64  `json:"Tail"`
}

type logMsg struct {
	Sender  string `json:"Sender"`
	Message string `json:"Message"`
}

func MakeLogMsg(sndr, msg string) *logMsg { // check size before sending
	return &logMsg{
		Sender:  sndr,
		Message: msg,
	}
}

// returns [<sender>]: <message>
func (m logMsg) String() string{
	return fmt.Sprintf("[%s]: %s", m.Sender, m.Message)
}

// creates a header which breifs the transmission receiver on what to expect from the incoming file.
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
