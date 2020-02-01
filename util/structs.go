package structs

type SaveFile struct {
	Files []FileHeader `json:"Files"`
}

type FileHeader struct {
	Filename  string `json:"Filename"`
	User      string `json:"User"`
	Date      string `json:"Date"`
	AuthToken string `json:"Authentication"`
	TailSize  int  `json:"Tail"`
	Blocks    int64  `json:"Blocks"`
}
