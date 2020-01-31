package structs

type SaveFile struct {
	Thoughts []ThoughtBody `json:"Thoughts"`
}

type ThoughtBody struct {
	Date    string `json:"Date"`
	User    string `json:"User"`
	Thought string `json:"Thought"`
}

type ThoughtHeader struct {
	AuthToken string `json:"Authentication"`
	Size      int    `json:"Size"`
}

type ThoughtBubble struct {
	Head ThoughtHeader `json:"Header"`
	Body ThoughtBody   `json:"Body"`
}
