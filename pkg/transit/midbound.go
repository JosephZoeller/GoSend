package transit

import (
	"io"
	"log"
	"net"
)

// PassHeader Reroutes a file header from lCon to sCon.
func PassHeader(lCon *net.Conn, sCon *net.Conn) (*fileHeader, error) {
	fHead, er := HeaderInbound(lCon)
	if er != nil {
		return nil, er
	}

	er = HeaderOutbound(fHead, sCon)
	if er != nil {
		return fHead, er
	}
	return fHead, nil
}

// PassFile Reroutes the file download stream from lCon to sCon.
// Requires a file header to determine the size and name of the file.
func PassFile(fHead *fileHeader, lCon *net.Conn, sCon *net.Conn) error {
	lc := *lCon
	sc := *sCon

	kb := fHead.Kilobytes
	tail := fHead.TailSize
	sendSize := int64(1024)

	for i := int64(0); i <= kb; i++ {
		if (i == kb) {
			sendSize = tail
		}
		_, er := io.CopyN(sc, lc, sendSize)
		//log.Println(n)
		if er != nil {
			return er
		}
	}
	log.Printf("[Pass File]: Successfully passed %s file.", fHead.Filename)
	return nil
}
