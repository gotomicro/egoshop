package user

import "github.com/i2eco/egoshop/appgo/model/trans"

type ReqList struct {
	Tid    int `form:"tid"`
	TypeId int `form:"typeId"`
	trans.ReqPage
}
