package trans

type ReqPage struct {
	Current  int    `json:"currentPage" form:"currentPage"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Sort     string `json:"sort" form:"sort"`
}
