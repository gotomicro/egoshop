package cart

type ReqInfo struct {
	ComSkuId int `form:"comSkuId"`
}

type ReqList struct {
	CartIds []int `json:"cartIds"`
	Current int   `form:"current"`
}

type ReqCreate struct {
	TypeId   int `json:"typeId"`
	ComSkuId int `json:"comSkuId"`
	Num      int `json:"num"`
	FeedId   int `json:"feedId"`
}

type ReqUpdate struct {
	TypeId   int `json:"typeId"`
	ComSkuId int `json:"comSkuId"`
	Quantity int `json:"quantity"`
}

type ReqDel struct {
	Ids []ReqDelOne `json:"ids"`
}

type ReqDelOne struct {
	TypeId    int `json:"typeId"`
	FeedId    int `json:"feedId"`
	ComSkuIds int `json:"comSkuIds"`
}

type ReqCheck struct {
	ComSkuIds []int `json:"comSkuIds"`
	IsCheck   int   `json:"isCheck"`
}
