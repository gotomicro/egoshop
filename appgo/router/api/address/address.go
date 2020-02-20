package address

import (
	"github.com/gin-gonic/gin"
	"github.com/i2eco/egoshop/appgo/dao"
	"github.com/i2eco/egoshop/appgo/model/mysql"
	"github.com/i2eco/egoshop/appgo/model/trans"
	"github.com/i2eco/egoshop/appgo/pkg/base"
	"github.com/i2eco/egoshop/appgo/pkg/code"
	"github.com/i2eco/egoshop/appgo/pkg/mus"
	"github.com/i2eco/egoshop/appgo/router/mdw"
)

func List(c *gin.Context) {
	uid := mdw.WechatUid(c)
	cnt, list := dao.Address.ListPage(c, mysql.Conds{
		"created_by": uid,
	}, &trans.ReqPage{
		Current:  1,
		PageSize: 100,
		Sort:     "id desc",
	})

	output := make([]mysql.Address, 0)
	for _, info := range list {
		newInfo := &info
		newInfo.WithTypeName()
		output = append(output, *newInfo)
	}
	base.JSONWechatList(c, output, cnt, 100)
}

func TypeList(c *gin.Context) {
	list, _ := dao.AddressType.List(c, mysql.Conds{})
	base.JSONOK(c, list)
}

func Info(c *gin.Context) {
	req := ReqInfo{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	uid := mdw.WechatUid(c)
	info, err := dao.Address.InfoX(c, mysql.Conds{
		"id":         req.Id,
		"created_by": uid,
	})
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	// 获取type name
	info.WithTypeName()
	base.JSON(c, code.MsgOk, info)
}

func Create(c *gin.Context) {
	req := ReqCreate{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	uid := mdw.WechatUid(c)

	var err error
	tx := mus.Db.Begin()
	if req.IsDefault == 1 {
		err = dao.Address.UpdateX(c, tx, mysql.Conds{
			"updated_by": uid,
		}, mysql.Ups{
			"is_default": 0,
		})
		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)

			return
		}
	}

	createAddress := &mysql.Address{
		Id:        0,
		CreatedBy: uid,
		UpdatedBy: uid,
		Name:      req.Name,
		Region:    req.Region,
		Detail:    req.Detail,
		TelPhone:  "",
		Mobile:    req.Mobile,
		ZipCode:   "",
		IsDefault: req.IsDefault,
		TypeId:    req.TypeId,
		StreetId:  0,
	}

	err = dao.Address.Create(c, tx, createAddress)

	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	tx.Commit()

	base.JSONOK(c)
}

func SetDefault(c *gin.Context) {
	req := ReqSetDefault{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	uid := mdw.WechatUid(c)
	err := dao.Address.UpdateX(c, mus.Db, mysql.Conds{
		"id":         req.Id,
		"updated_by": uid,
	}, mysql.Ups{
		"is_default": 1,
	})
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	base.JSONOK(c)
}

func Default(c *gin.Context) {
	uid := mdw.WechatUid(c)
	info, _ := dao.Address.InfoX(c, mysql.Conds{
		"created_by": uid,
		"is_default": 1,
	})
	base.JSONOK(c, info)
}

func Update(c *gin.Context) {
	req := ReqUpdate{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	uid := mdw.WechatUid(c)

	var err error
	tx := mus.Db.Begin()
	if req.IsDefault == 1 {
		err = dao.Address.UpdateX(c, tx, mysql.Conds{
			"updated_by": uid,
		}, mysql.Ups{
			"is_default": 0,
		})
		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)
			return
		}
	}

	err = dao.Address.UpdateX(c, tx, mysql.Conds{
		"id": req.Id,
	}, mysql.Ups{
		"name":       req.Name,
		"region":     req.Region,
		"detail":     req.Detail,
		"tel_phone":  "",
		"mobile":     req.Mobile,
		"is_default": req.IsDefault,
		"type_id":    req.TypeId,
	})

	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	tx.Commit()

	base.JSONOK(c)
}

func Del(c *gin.Context) {
	req := ReqDel{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	uid := mdw.WechatUid(c)
	err := dao.Address.DeleteX(c, mus.Db, mysql.Conds{
		"created_by": uid,
		"id":         req.Id,
	})

	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	base.JSONOK(c)

}
