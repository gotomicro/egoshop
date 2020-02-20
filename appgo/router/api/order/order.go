package order

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/constx"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
	"github.com/goecology/egoshop/appgo/pkg/mus"
	"github.com/goecology/egoshop/appgo/router/mdw"
	"github.com/jinzhu/gorm"
	"github.com/satori/uuid"
)

func StateNum(c *gin.Context) {
	uid := mdw.WechatUid(c)
	base.JSON(c, code.MsgOk, map[string]interface{}{
		"stateNew":     dao.Order.GetStateCnt(uid, constx.OrderStateNew),
		"stateSend":    dao.Order.GetStateCnt(uid, constx.OrderStateSend),
		"stateSuccess": dao.Order.GetStateCnt(uid, constx.OrderStateSuccess),
		"stateClose":   dao.Order.GetStateCnt(uid, constx.OrderStateClose),
	})
}

/**
 * 创建订单
 */
func Create(c *gin.Context) {
	req := ReqCreate{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request error")
		return
	}
	uid := mdw.WechatUid(c)

	addressInfo, err := dao.Address.InfoX(c, mysql.Conds{
		"id":         req.AddressId,
		"created_by": uid,
	})
	if err != nil {
		err = errors.New("收货地址没找到")
		return
	}

	reqPage := trans.ReqPage{
		Current:  1,
		PageSize: 10000,
		Sort:     "updated_at desc",
	}

	list, _ := dao.Cart.ListAddition(c, uid, req.CartIds, &reqPage)

	tx := mus.Db.Begin()

	create := &mysql.OrderPay{
		PaySn:    strings.ReplaceAll(uuid.NewV4().String(), "-", ""),
		Uid:      uid,
		PayState: constx.OrderPayStateNew,
	}

	err = dao.OrderPay.Create(c, tx, create)
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	var calculateValue BuyCalculate
	calculateValue, err = calculate(c, list)
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	orderCreate := &mysql.Order{
		Id:               0,
		Sn:               strings.ReplaceAll(uuid.NewV4().String(), "-", ""),
		PaySn:            create.PaySn,
		Uid:              uid,
		CreatedBy:        uid,
		UpdatedBy:        uid,
		UserName:         mdw.WechatUserName(c),
		UserPhone:        addressInfo.Mobile,
		UserEmail:        "",
		PaymentCode:      "online",
		PayName:          "online",
		PaymentTime:      0,
		FinishedTime:     0,
		GoodsAmount:      calculateValue.ComAmount,
		GoodsNum:         calculateValue.ComNum,
		Amount:           calculateValue.ComAmount, //目前没有运输费,所以两个加个相等
		PdAmount:         0,
		State:            constx.OrderStateNew,
		RefundAmount:     0,
		RefundState:      0,
		LockState:        0,
		TradeNo:          "",
		OutRequestNo:     "",
		From:             uid,
		ReviseAmount:     0,
		ReviseFreightFee: 0,
	}

	err = dao.Order.Create(c, tx, orderCreate)
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	// 拓展订单表创建
	orderExtendCreate := &mysql.OrderExtend{
		Id: orderCreate.Id,
		ReceiverInfo: mysql.OrderExtendReceiverInfoJson{
			Name:   addressInfo.Name,
			Detail: addressInfo.Detail,
			Phone:  addressInfo.Mobile,
			//Type:          addressInfo.TypeId,
			Address: addressInfo.Region,
		},
		ReceiverName:       addressInfo.Name,
		ReceiverPhone:      addressInfo.Mobile,
		ReceiverProvinceId: addressInfo.ProvinceId,
		ReceiverCityId:     addressInfo.CityId,
		ReceiverAreaId:     addressInfo.AreaId,
	}

	err = dao.OrderExtend.Create(c, tx, orderExtendCreate)
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	cartIds := make([]int, 0)
	for _, value := range list {
		cartIds = append(cartIds, value.Id)
		var create4 *mysql.OrderGoods
		if value.TypeId == 1 {
			create4 = &mysql.OrderGoods{
				Uid:      mdw.WechatUid(c),
				OrderId:  orderCreate.Id,
				TypeId:   1,
				FeedId:   value.FeedId,
				Title:    value.Title,
				Price:    value.Price,
				PayPrice: float64(value.Num) * value.Price,
				Num:      value.Num,
				Cover:    value.Cover,
			}
		} else if value.TypeId == 2 {
			create4 = &mysql.OrderGoods{
				Uid:      mdw.WechatUid(c),
				OrderId:  orderCreate.Id,
				TypeId:   2,
				ComId:    value.ComId,
				ComSkuId: value.ComSkuId,
				Title:    value.Title,
				Price:    value.Price,
				PayPrice: float64(value.Num) * value.Price,
				Num:      value.Num,
				Cover:    value.Cover,
				ComType:  1,
			}
		}

		err = dao.OrderGoods.Create(c, tx, create4)
		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)
			return
		}
	}

	createLog := &mysql.OrderLog{
		OrderId:    orderCreate.Id,
		Msg:        "买家下单",
		Role:       "buyer",
		OrderState: constx.OrderStateNew,
	}

	err = dao.OrderLog.Create(c, tx, createLog)
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	// todo
	err = updateStorageNum(tx, list)
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	// todo 优化,delete by in
	for _, cartId := range cartIds {
		err = dao.Cart.DeleteX(c, tx, mysql.Conds{
			"id": cartId,
		})

		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)
			return
		}
	}

	tx.Commit()
	base.JSONOK(c, gin.H{
		"paySn":   create.PaySn,
		"orderId": orderCreate.Id,
	})
}

func Info(c *gin.Context) {
	req := ReqInfo{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr)
		return
	}
	uid := mdw.WechatUid(c)

	orderInfo, err := dao.Order.GetOrderInfo(c, req.Id, uid)

	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	base.JSON(c, code.MsgOk, orderInfo)
}

func Delete(c *gin.Context) {
	req := ReqDelete{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr)
		return
	}

	uid := mdw.WechatUid(c)
	orderInfo, err := dao.Order.InfoX(c, mysql.Conds{
		"id":         req.Id,
		"created_by": uid,
	})
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	isCancel := dao.Order.GetOrderOperateState("user_cancel", orderInfo)
	// 如果不可以取消，那么报错
	if isCancel == false {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	err = dao.Order.Delete(c, mus.Db, req.Id)
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	base.JSON(c, code.MsgOk)
}

/**
id;
    sn;
    pay_sn;
    create_time;
    payment_code;
    pay_name;
    payment_time;
    finnshed_time;
    goods_amount;
    goods_num;
    amount;
    pd_amount;
    freight_fee;
    freight_unified_fee;
    freight_template_fee;
    state;
    refund_amount;
    refund_state;
    lock_state;
    delay_time;
    tracking_no;
    evaluate_state;
    trade_no;
    state_desc;
    payment_name;
    extend_order_extend;
    extend_order_goods;

    if_cancel;
    if_pay;
    if_evaluate;
    if_receive;
*/

func List(c *gin.Context) {
	req := ReqOrderListInfo{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}
	uid := mdw.WechatUid(c)
	conds := mysql.Conds{
		"uid": uid,
	}
	req.ReqPage.PageSize = 10

	// 订单状态查询
	if req.StateType != "" && req.StateType != "all" {
		state, ok := constx.OrderStates[req.StateType]
		if !ok {
			base.JSON(c, code.MsgErr, "request app list params is error")
			return
		}
		conds["state"] = state
	}

	total, orderMap := dao.Order.GetOrderFullsPage(c, &req.ReqPage, conds, map[string]bool{
		"order_goods":  true,
		"order_extend": true,
	})

	for _, value := range orderMap {
		// 显示取消订单
		value.IsCancel = dao.Order.GetOrderOperateState("user_cancel", value.Order)
		// 显示是否需能支付（todo 计算后台过期时间）
		value.IsPay = dao.Order.GetOrderOperateState("user_pay", value.Order)
		// 显示退款取消订单
		value.IsRefundCancel = dao.Order.GetOrderOperateState("refund_cancel", value.Order)
		// 显示投诉
		value.IsComplain = dao.Order.GetOrderOperateState("complain", value.Order)
		// 显示收货
		value.IsReceive = dao.Order.GetOrderOperateState("receive", value.Order)
		// 显示锁定中
		value.IsLock = dao.Order.GetOrderOperateState("lock", value.Order)
		// 显示物流跟踪
		value.IsDeliver = dao.Order.GetOrderOperateState("deliver", value.Order)
		// 显示评价
		value.IsEvaluate = dao.Order.GetOrderOperateState("evaluate", value.Order)
	}

	output := make([]dao.OrderFull, 0)
	for _, value := range orderMap {
		output = append(output, *value)
	}
	// todo要按时间排序		Sort:     "create_time desc",
	base.JSONWechatList(c, output, total, req.ReqPage.PageSize)
}

//func Cancel(c *gin.Context) {
//	req := trans.ReqOrderCancel{}
//	if err := c.Bind(&req); err != nil {
//		base.JSON(c, base.MsgErr, "request app list params is error")
//		return
//	}
//	uid := mdw.WechatUid(c)
//	openId := mdw.OpenId(c)
//
//	orderInfo, err := dao.Order.InfoByIdAndUid(req.Id, uid, openId)
//	if err != nil {
//		base.JSON(c, base.MsgErr, "info not exist")
//		return
//	}
//	fmt.Println(orderInfo)
//
//}
//
///**
//id;
//    sn;
//    pay_sn;
//    create_time;
//    payment_code;
//    pay_name;
//    payment_time;
//    finnshed_time;
//    com_amount;
//    com_num;
//    amount;
//    pd_amount;
//    freight_fee;
//    freight_unified_fee;
//    freight_template_fee;
//    state;
//    refund_amount;
//    refund_state;
//    lock_state;
//    delay_time;
//    tracking_no;
//    evaluate_state;
//    trade_no;
//    state_desc;
//    payment_name;
//    extend_order_extend;
//    extend_order_com;
//
//    if_cancel;
//    if_pay;
//    if_evaluate;
//    if_receive;
//*/
//
//func List(c *gin.Context) {
//	req := trans.ReqOrderListInfo{}
//	if err := c.Bind(&req); err != nil {
//		base.JSON(c, base.MsgErr, "request is error")
//		return
//	}
//	uid := wechat.WechatUid(c)
//	openId := mdw.OpenId(c)
//	conds := mysql.Conds{
//		"uid":     uid,
//		"open_id": openId,
//	}
//
//	// 订单状态查询
//	if req.StateType != "" {
//		state, ok := common.OrderStates[req.StateType]
//		if !ok {
//			base.JSON(c, base.MsgErr, "request app list params is error")
//			return
//		}
//		conds["state"] = state
//	}
//
//	total, orderMap := dao.Order.GetOrderFullsPage(c, trans.ReqPage{
//		CurrentPage:  req.Page - 1,
//		PageSize: req.Rows,
//	}, conds, map[string]bool{
//		"order_com":    true,
//		"order_extend": true,
//	})
//
//	for _, value := range orderMap {
//		// 显示取消订单
//		value.IsCancel = dao.Order.GetOrderOperateState("user_cancel", value.Order)
//		// 显示是否需能支付（todo 计算后台过期时间）
//		value.IsPay = dao.Order.GetOrderOperateState("user_pay", value.Order)
//		// 显示退款取消订单
//		value.IsRefundCancel = dao.Order.GetOrderOperateState("refund_cancel", value.Order)
//		// 显示投诉
//		value.IsComplain = dao.Order.GetOrderOperateState("complain", value.Order)
//		// 显示收货
//		value.IsReceive = dao.Order.GetOrderOperateState("receive", value.Order)
//		// 显示锁定中
//		value.IsLock = dao.Order.GetOrderOperateState("lock", value.Order)
//		// 显示物流跟踪
//		value.IsDeliver = dao.Order.GetOrderOperateState("deliver", value.Order)
//		// 显示评价
//		value.IsEvaluate = dao.Order.GetOrderOperateState("evaluate", value.Order)
//	}
//
//	output := make([]dao.OrderFull, 0)
//	for _, value := range orderMap {
//		output = append(output, *value)
//	}
//	// todo要按时间排序		Sort:     "create_time desc",
//
//	base.JSONList(c, output, total)
//}
//
//func ConfirmReceipt(c *gin.Context) {
//	req := trans.ReqOrderConfirmReceipt{}
//	if err := c.Bind(&req); err != nil {
//		base.JSON(c, base.MsgErr, "request app list params is error")
//		return
//	}
//	uid := wechat.WechatUid(c)
//	openId := mdw.OpenId(c)
//
//	orderInfo, err := dao.Order.InfoByIdAndUid(req.Id, uid, openId)
//	if err != nil {
//		base.JSON(c, base.MsgErr, "info not exist")
//		return
//	}
//
//	flag := dao.Order.UserChangeState(c, "order_receive", orderInfo, uid, wechat.WechatUserName(c), req.StateRemark)
//	if !flag {
//		base.JSON(c, base.MsgErr, "change state error")
//		return
//	}
//	base.JSON(c, base.MsgOk, "change state ok")
//	return
//
//}
//
//func Logistics(c *gin.Context) {
//	req := trans.ReqOrderLogistics{}
//	if err := c.Bind(&req); err != nil {
//		base.JSON(c, base.MsgErr, "request app list params is error")
//		return
//	}
//
//	_, err := dao.Order.InfoX(c, mysql.Conds{
//		"id":    req.Id,
//		"state": mysql.Cond{">", 30},
//		"uid":   wechat.WechatUid(c),
//	})
//	if err != nil {
//		base.JSON(c, base.MsgErr, "request app list params is error")
//		return
//	}
//
//	orderExtendInfo, err := dao.OrderExtend.InfoX(c, mysql.Conds{
//		"id": req.Id,
//	})
//
//	if err != nil {
//		base.JSON(c, base.MsgErr, "request app list params is error")
//		return
//	}
//
//	if orderExtendInfo.TrackingNo != "" && orderExtendInfo.ShipperId > 0 && orderExtendInfo.ExpressId > 0 {
//		info, err := dao.Express.InfoX(c, mysql.Conds{
//			"id": orderExtendInfo.ExpressId,
//		})
//		if err != nil {
//			base.JSON(c, base.MsgOk, "ok", map[string]interface{}{
//				"url": "",
//			})
//			return
//		}
//
//		base.JSON(c, base.MsgOk, "ok", map[string]interface{}{
//			"url": "https://m.kuaidi100.com/index_all.html?type=" + info.Kuaidi100Code + "&postid=" + orderExtendInfo.TrackingNo,
//		})
//		return
//	}
//	base.JSON(c, base.MsgOk, "ok", map[string]interface{}{
//		"url": "",
//	})
//}
//
//func ComList(c *gin.Context) {
//	req := trans.ReqOrderComList{}
//	if err := c.Bind(&req); err != nil {
//		base.JSON(c, base.MsgErr, "request app list params is error")
//		return
//	}
//
//	cnt, list := dao.OrderGoods.ListPage(c, mysql.Conds{
//		"order_id": req.Id,
//		"uid":      wechat.WechatUid(c),
//	}, trans.ReqPage{
//		CurrentPage:  0,
//		PageSize: 10000,
//		Sort:     "id asc",
//	})
//	base.JSONList(c, list, cnt)
//}
//
//func ComInfo(c *gin.Context) {
//	req := trans.ReqOrderComInfo{}
//	if err := c.Bind(&req); err != nil {
//		base.JSON(c, base.MsgErr, "request app list params is error")
//		return
//	}
//
//	resp, _ := dao.OrderGoods.InfoX(c, mysql.Conds{
//		"order_id": req.Id,
//		"uid":      wechat.WechatUid(c),
//	})
//	base.JSON(c, base.MsgOk, map[string]interface{}{
//		"info": resp,
//	})
//}

type BuyCalculate struct {
	ComAmount float64
	ComNum    int
}

func calculate(c *gin.Context, list []mysql.Cart) (output BuyCalculate, err error) {
	var comAmount float64
	var comNum int

	for _, value := range list {
		if value.Num > value.Stock {
			err = errors.New("库存不够")
			return
		}

		comAmount += value.Price * float64(value.Num)
		comNum += value.Num
	}

	output = BuyCalculate{
		ComAmount: comAmount,
		ComNum:    comNum,
	}
	return
}

func updateStorageNum(tx *gorm.DB, list []mysql.Cart) (err error) {

	comSkuInfos := make(map[int]comSkuInfo, 0)

	for _, value := range list {
		// 同一款产品的sku数据相加更新到主表
		if value.TypeId == constx.TypeCom {
			comId := value.ComId
			if tmpValue, ok := comSkuInfos[comId]; !ok {
				comSkuInfos[comId] = comSkuInfo{
					Id:      value.ComId,
					Stock:   0,
					SaleNum: value.Num,
				}
			} else {
				tmpValue.SaleNum += value.Num
				comSkuInfos[comId] = tmpValue
			}
			// 更新sku商品里的统计值
			err = tx.Exec("UPDATE com_sku set `stock`=stock-? , `sale_num`=sale_num+? WHERE id=?", value.Num, value.Num, value.ComSkuId).Error
			if err != nil {
				return
			}
		}

	}
	// 更新主表sku商品里的统计值
	for _, value := range comSkuInfos {
		err = tx.Exec("UPDATE com set `stock`=stock-? , `sale_num`=sale_num+? WHERE id=?", value.SaleNum, value.SaleNum, value.Id).Error
		if err != nil {
			return
		}
	}
	return
}

type comSkuInfo struct {
	Id      int
	Stock   int
	SaleNum int
}
