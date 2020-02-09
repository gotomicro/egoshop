package admineditor

type ReqContentSave struct {
	Id          int    `json:"id"`
	Type        int    `json:"type"`
	Name        string `json:"name"`
	PreMarkdown string `json:"markdown"`
	PreHtml     string `json:"html"`
}

type ReqUpload struct {
	Id    int    `json:"id"`
	Type  int    `json:"type"`
	Space string `json:"space"`
	Image string `json:"image"`
}

type ReqRelease struct {
	Id   int `json:"id"`
	Type int `json:"type"`
}
