package transit

type saveFile struct {
	Files []fileHeader `json:"Files"`
}

type fileHeader struct {
	Filename  string `json:"Filename"`
	User      string `json:"User"`
	Date      string `json:"Date"`
	AuthToken string `json:"Authentication"`
	Blocks    int64  `json:"Blocks"`
	TailSize  int64    `json:"Tail"`
}
