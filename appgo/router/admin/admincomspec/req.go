package admincomspec

import (
	"github.com/goecology/egoshop/appgo/model/mysql"
)

type RespComSpecList struct {
	Id     int                  `json:"id"`
	Name   string               `json:"name"`
	Values []mysql.ComSpecValue `json:"values"`
}

type ReqComSpecCreate struct {
	Name string `json:"name"`
}

type ReqComSpecValueCreate struct {
	SpecId int    `json:"specId"`
	Name   string `json:"name"`
}
