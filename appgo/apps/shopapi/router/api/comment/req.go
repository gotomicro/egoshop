package comment

type ReqSubmit struct {
	Gid     int    `json:"gid"`
	TypeId  int    `json:"typeId"`
	Content string `json:"content"`
}

type ReqList struct {
	CurrentPage int `form:"currentPage"`
	Gid         int `form:"gid"`
	TypeId      int `form:"typeId"`
}
