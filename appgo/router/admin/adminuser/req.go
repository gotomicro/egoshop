package adminuser

import "github.com/goecology/egoshop/appgo/model/trans"

type ReqList struct {
	Name string `form:"name"`
	trans.ReqPage
}
