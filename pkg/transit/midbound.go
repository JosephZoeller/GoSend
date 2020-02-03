package transit

import (
	"io"
	"net"

	"github.com/JosephZoeller/gmg/pkg/logUtil"
)

// PassHeader Reroutes a file header from lCon to sCon.
func PassHeader(lCon *net.Conn, sCon *net.Conn) (*fileHeader, error) {
	fHead, er := HeaderInbound(lCon)
	if er != nil {
		return nil, logUtil.FormatError("Transit PassHeader", er)
	}

	er = HeaderOutbound(fHead, sCon)
	if er != nil {
		return fHead, logUtil.FormatError("Transit PassHeader", er)
	}
	return fHead, nil
}

// PassFile Reroutes the file download stream from lCon to sCon.
// Requires a file header to determine the size and name of the file.
func PassFile(fHead *fileHeader, lCon *net.Conn, sCon *net.Conn) error {
	lc := *lCon
	sc := *sCon

	for i := int64(0); i <= fHead.Blocks; i++ {
		_, er := io.CopyN(sc, lc, 1024)
		if er != nil {
			return logUtil.FormatError("Transit PassFile", er)
		}
	}

	return nil
}
