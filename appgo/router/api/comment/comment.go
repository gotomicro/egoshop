package comment

import (
	"errors"

	"github.com/i2eco/egoshop/appgo/service"

	"github.com/gin-gonic/gin"
	"github.com/i2eco/egoshop/appgo/dao"
	"github.com/i2eco/egoshop/appgo/model/common"
	"github.com/i2eco/egoshop/appgo/model/constx"
	"github.com/i2eco/egoshop/appgo/model/mysql"
	"github.com/i2eco/egoshop/appgo/model/trans"
	"github.com/i2eco/egoshop/appgo/pkg/base"
	"github.com/i2eco/egoshop/appgo/pkg/code"
	"github.com/i2eco/egoshop/appgo/pkg/mus"
	"github.com/i2eco/egoshop/appgo/router/mdw"
)

// 写评论
func Create(c *gin.Context) {
	reqParam := ReqSubmit{}
	var err error
	err = c.Bind(&reqParam)
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	// todo合并验证解析
	goodsInfo, err := mdw.ValidateGoodsByParam(common.Goods{
		Gid:    reqParam.Gid,
		TypeId: reqParam.TypeId,
	})
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	// 评论内容
	// todo 验证大小
	if reqParam.Content == "" {
		base.JSONErr(c, code.MsgErr, errors.New("content is empty"))
		return
	}

	uid := mdw.WechatUid(c)
	resp, err := dao.UserOpen.InfoX(c, mysql.Conds{
		"uid": uid,
	})
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	create := &mysql.Comment{
		GoodsId:  goodsInfo.Gid,
		TypeId:   goodsInfo.TypeId,
		Content:  reqParam.Content,
		Uid:      uid,
		Score:    0, //todo 评分均为0
		Nickname: resp.Nickname,
		Avatar:   resp.Avatar,
	}

	err = dao.Comment.Create(c, mus.Db, create)
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	// 增加积分
	service.QueuePoint.Push(service.PointTask{
		TypeData: constx.PointComment,
		Uid:      uid,
	})
	base.JSONOK(c)
}

// 获取评论
func List(c *gin.Context) {
	req := ReqList{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	reqPage := trans.ReqPage{
		Current:  req.CurrentPage,
		PageSize: 10,
		Sort:     "id desc",
	}
	// todo 评论分页
	total, resp := dao.Comment.ListPage(c, mysql.Conds{"goods_id": req.Gid, "type_id": req.TypeId}, &reqPage)
	base.JSONWechatList(c, resp, total, reqPage.PageSize)
}

func ListTop3(c *gin.Context) {
	req := ReqList{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	reqPage := trans.ReqPage{
		Current:  1,
		PageSize: 3,
		Sort:     "id desc",
	}
	// todo 评论分页
	total, resp := dao.Comment.ListPage(c, mysql.Conds{"goods_id": req.Gid, "type_id": req.TypeId}, &reqPage)
	base.JSONWechatList(c, resp, total, reqPage.PageSize)
}
