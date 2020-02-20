package admincomspec

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
	"github.com/goecology/egoshop/appgo/pkg/mus"
)

func List(c *gin.Context) {
	total, list := dao.ComSpec.ListPage(c, mysql.Conds{}, &trans.ReqPage{
		PageSize: 10000,
		Sort:     "sort desc",
	})
	resp := make([]RespComSpecList, 0)
	var wg sync.WaitGroup
	wg.Add(len(list))
	for _, value := range list {
		go func(value mysql.ComSpec) {
			_, list := dao.ComSpecValue.ListPage(c, mysql.Conds{
				"spec_id": value.Id,
			}, &trans.ReqPage{
				PageSize: 10000,
				Sort:     "sort desc",
			})
			resp = append(resp, RespComSpecList{
				Id:     value.Id,
				Name:   value.Name,
				Values: list,
			})
			wg.Done()
		}(value)
	}
	wg.Wait()
	base.JSONList(c, resp, 0, 0, total)
}

func Create(c *gin.Context) {
	req := ReqComSpecCreate{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	err := dao.ComSpec.Create(c, mus.Db, &mysql.ComSpec{
		Name: req.Name,
	})
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	base.JSONOK(c)
}

func ValueCreate(c *gin.Context) {
	req := ReqComSpecValueCreate{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	err := dao.ComSpecValue.Create(c, mus.Db, &mysql.ComSpecValue{
		SpecId: req.SpecId,
		Name:   req.Name,
		Sort:   0,
		Color:  "",
		Img:    "",
	})
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	base.JSONOK(c)
}
