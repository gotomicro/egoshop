package adminuser

import "github.com/i2eco/egoshop/appgo/model/trans"

type ReqList struct {
	Name string `form:"name"`
	trans.ReqPage
}
