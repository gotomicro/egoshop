package cate

type ReqList struct {
	CateId    int  `json:"cateid"`
	CateChild bool `json:"catechild"` //是否获取子类
}
