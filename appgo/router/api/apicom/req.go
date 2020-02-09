package apicom

type ReqList struct {
	//Cids []int `form:"cids"`
	CurrentPage int    `form:"currentPage"`
	IsNew       int    `form:"isNew"`
	Order       string `form:"order"`
}
