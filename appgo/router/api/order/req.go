package order

import "github.com/goecology/egoshop/appgo/model/trans"

type ReqInfo struct {
	Id int `form:"id"`
}

type ReqDelete struct {
	Id int `form:"id"`
}

type ReqOrderCancel struct {
	Id          int `json:"id"`
	StateRemark int `json:"state_remark"`
}

type ReqOrderListInfo struct {
	trans.ReqPage
	StateType string `form:"stateType"`
}

type ReqOrderConfirmReceipt struct {
	Id          int    `json:"id"`
	StateRemark string `json:"state_remark"`
}

type ReqOrderGoodsInfo struct {
	Id int `json:"id"`
}

type ReqOrderLogistics struct {
	Id int `json:"id"`
}

type ReqCreate struct {
	CartIds   []int  `json:"cartIds"`
	AddressId int    `json:"addressId"`
	Remark    string `json:"remark"`
}
