package admincom

import (
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
)

type ReqList struct {
	Name    string `form:"name"`
	Author  string `form:"author"`
	Address string `form:"address"`
	trans.ReqPage
}

type ReqCreateOrUpdate struct {
	Id         int            `json:"id"`
	Cid        int            `json:"cid"`
	Title      string         `json:"title"`
	SubTitle   string         `json:"subTitle"`
	Gallery    []string       `json:"gallery"`
	SaleTime   string         `json:"saleTime"`
	Cids       []string       `json:"cids"`
	SkuList    []mysql.ComSku `json:"skuList"`
	FreightFee float64        `json:"freightFee"` // 运费
	FreightId  int            `json:"freightId"`  // 运费模板id
}

type ReqUpdate struct {
	Id       int      `json:"id"`
	Cid      int      `json:"cid"`
	Title    string   `json:"title"`
	SubTitle string   `json:"subTitle"`
	Gallery  []string `json:"gallery"`
	SaleTime string   `json:"saleTime"`
}

type ReqRemove struct {
	Id int `json:"id"`
}

type ReqOnSale struct {
	Ids []int `json:"ids"`
}
