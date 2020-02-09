package transit

import (
	"fmt"
	"os"
)

type fileHeader struct {
	Filename  string `json:"Filename"`
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
func (m logMsg) String() string {
	return fmt.Sprintf("[%s]: %s", m.Sender, m.Message)
}

// creates a header which breifs the transmission receiver on what to expect from the incoming file.
func MakeHeader(file *os.File) (*fileHeader, error) {
	fstat, er := file.Stat()
	if er != nil {
		return &fileHeader{}, er
	}

	return &fileHeader{
		Kilobytes: fstat.Size() / 1024,
		TailSize:  fstat.Size() % 1024,
		Filename:  fstat.Name(),
	}, nil
}
