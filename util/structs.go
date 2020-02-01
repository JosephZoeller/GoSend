package structs

type SaveFile struct {
	Files []FileHeader `json:"Files"`
}

type FileHeader struct {
	Filename  string `json:"Filename"`
	User      string `json:"User"`
	Date      string `json:"Date"`
	AuthToken string `json:"Authentication"`
	Blocks    int64  `json:"Blocks"`
	TailSize  int    `json:"Tail"`
}
