package admincomcate

import "github.com/goecology/egoshop/appgo/model/trans"

type ReqList struct {
	Name    string `form:"name"`
	Author  string `form:"author"`
	Address string `form:"address"`
	trans.ReqPage
}

type ReqInfo struct {
	Id int `form:"id"`
}

type ReqCreate struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Pid  int    `json:"pid"`
}

type ReqUpdate struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type ReqRemove struct {
	Id int `json:"id"`
}
