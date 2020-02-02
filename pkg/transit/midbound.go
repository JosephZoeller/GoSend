package transit

import (
	"io"
	"net"
)

func PassHeader(lCon *net.Conn, sCon *net.Conn) (*fileHeader, error) {
	fHead, er := HeaderInbound(lCon)
	if er != nil {
		return nil, er
	}
	return fHead, HeaderOutbound(fHead, sCon)
}

func PassFile(fHead *fileHeader, lCon *net.Conn, sCon *net.Conn) error {
	lc := *lCon
	sc := *sCon

	for i := int64(0); i <= fHead.Blocks; i++ {
		_, er := io.CopyN(sc, lc, 1024)
		if er != nil {
			return er
		}
	}

	return nil
}
