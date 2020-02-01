package buy

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/apps/shopapi/pkg/conf"
	"github.com/goecology/egoshop/appgo/apps/shopapi/pkg/mus"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router/mdw"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/constx"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
	"github.com/jinzhu/gorm"
	"github.com/milkbobo/gopay"
	"github.com/milkbobo/gopay/common"
	"github.com/milkbobo/gopay/constant"
	"go.uber.org/zap"
)

type Req struct {
	OrderType      string `json:"orderType"`
	PaySn          string `json:"paySn"`
	PaymentCode    string `json:"paymentCode"`
	PaymentChannel string `json:"paymentChannel"`
}

func Calculate(c *gin.Context) {
	req := ReqBuyCalculate{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgParamErr, err)
		return
	}

	uid := mdw.WechatUid(c)

	addressInfo, err := dao.Address.InfoX(c, mysql.Conds{
		"id":         req.AddressId,
		"created_by": uid,
	})
	if err != nil && err == gorm.ErrRecordNotFound {
		base.JSONErr(c, code.PayCalculateNotFoundAddress, err)
		return
	}

	reqPage := trans.ReqPage{
		Current:  1,
		PageSize: 10000,
		Sort:     "updated_at desc",
	}

	list, _ := dao.Cart.ListAddition(c, uid, req.CartIds, &reqPage)

	output, err := calculate(addressInfo, list)
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	base.JSON(c, code.MsgOk, output)
}

// Pay 订单支付接口
func Pay(c *gin.Context) {
	uid := mdw.WechatUid(c)

	req := Req{}
	// 验证参数是否合法
	err := c.Bind(&req)
	if err != nil {
		base.JSONErr(c, code.PayParamErr, err)
		return
	}

	order, err := dao.Order.InfoX(c, mysql.Conds{
		"pay_sn": req.PaySn,
		"uid":    uid,
		"state":  constx.OrderStateNew,
	})
	if err != nil && err == gorm.ErrRecordNotFound {
		err = errors.New("该订单不存在")
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	//_, _, err, codeStatus := payCheck(c, uid)
	//if err != nil {
	//	base.JSONErr(c, codeStatus, err)
	//	return
	//}

	_, err = dao.OrderPay.InfoX(c, mysql.Conds{
		"pay_sn":    req.PaySn,
		"pay_state": constx.OrderPayStateNew,
	})

	if err != nil && err == gorm.ErrRecordNotFound {
		err = errors.New("该订单不存在")
		return
	}

	payAmount := order.Amount
	if order.ReviseAmount > 0 {
		payAmount = order.ReviseAmount
	}

	result, err, codeStatus := payResult(c, req, uid, payAmount, order.PaySn)
	if err != nil {
		base.JSONErr(c, codeStatus, err)
		return
	}
	base.JSONOK(c, result)
	return
}

// payResult 返回支付json结果
func payResult(c *gin.Context, req Req, uid int, payAmount float64, paySn string) (result map[string]string, err error, codeStatus int) {
	result = make(map[string]string)
	switch req.PaymentChannel {
	case "wechat_mini":
		var userOpen mysql.UserOpen
		userOpen, err = dao.UserOpen.InfoX(c, mysql.Conds{"uid": uid, "genre": dao.GenreWechatMini})
		// 暂未获取到微信关联用户信息
		if err != nil {
			codeStatus = code.PayUserOpenErr
			return
		}
		charge := new(common.Charge)
		charge.PayMethod = constant.WECHAT_MINI_PROGRAM
		charge.MoneyFee = payAmount
		charge.TradeNum = paySn
		charge.Describe = fmt.Sprintf("商品购买_%s", paySn)
		charge.CallbackURL = conf.Conf.App.WechatPay.CallbackApi
		charge.OpenID = userOpen.MiniOpenid
		result, err = gopay.Pay(charge)
		if err != nil {
			codeStatus = code.PayWechatPayErr
			return
		}
		result["payAmount"] = strconv.FormatFloat(payAmount, 'E', -1, 32)
		break
	}
	return
}

func WXPayedNotify(c *gin.Context) {
	// 返回支付结果
	wechatResult, err := gopay.WeChatAppCallback(c.Writer, c.Request)
	mus.Logger.Info("[buy] WeChatAppCallback", zap.Any("wechatResult", wechatResult))
	if err != nil {
		base.JSONErr(c, code.PayPayedErr, err)
		return
	}

	// 检查返回值
	if wechatResult.ResultCode != "SUCCESS" {
		c.JSON(200, gin.H{"code": code.PayPayedErr, "msg": ""})
		return
	}

	// 接下来处理自己的逻辑
	err = dao.Order.Pay(c, wechatResult.OutTradeNO, "wechat", wechatResult.TransactionID)
	if err != nil {
		base.JSONErr(c, code.PayPayedErr, err)
		return
	}
	return
}

func calculate(addressInfo mysql.Address, list []mysql.Cart) (output RespBuyCalculate, err error) {
	freightList := make([]Freight, 0)

	// 根据运费计算规则一
	var freightUnifiedFee float64
	var freightTemplateFee float64
	var goodsAmount float64
	var payFreightFee float64
	var goodsNum int

	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	base.JSON(c, base.MsgErr, "数据库存在异常")
	// 	return
	// }

	for _, value := range list {
		if value.Num > value.Stock {
			err = errors.New("库存不够")
			return
		}
		valuePointer := &value
		var cost float64
		cost, err = valuePointer.FreightFeeByAddress(addressInfo)
		if err != nil {
			err = errors.New("数据异常")
			return
		}

		freightWay := valuePointer.GetFreightWay()

		freightList = append(freightList, Freight{
			ComSkuId:   valuePointer.ComSkuId,
			FreightFee: float64(cost),
			FreightWay: freightWay,
		})

		if freightWay == "goods_freight_unified" {
			// todo 为什么是=
			freightUnifiedFee = valuePointer.ComFreightFee
		} else {
			freightTemplateFee += valuePointer.ComFreightFee
		}

		goodsAmount += valuePointer.Price * float64(valuePointer.Num)
		goodsNum += valuePointer.Num
	}

	payFreightFee = freightUnifiedFee + freightTemplateFee

	output = RespBuyCalculate{
		ComAmount:          goodsAmount,
		PayAmount:          goodsAmount + payFreightFee,
		ComFreightList:     freightList,
		FreightUnifiedFee:  freightUnifiedFee,
		FreightTemplateFee: freightTemplateFee,
		PayFreightFee:      payFreightFee,
		SubTotal:           goodsAmount + payFreightFee,
		ComNum:             goodsNum,
	}
	return
}
