package structs

type SaveFile struct {
	Files []FileHeader `json:"Thoughts"`
}

type FileHeader struct {
	Filename  string `json:"Filename"`
	User      string `json:"User"`
	Date      string `json:"Date"`
	AuthToken string `json:"Authentication"`
	TailSize  int64  `json:"Tail"`
	Blocks    int64  `json:"Blocks"`
}
