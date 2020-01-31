package dao

import (
	"errors"

	"github.com/goecology/egoshop/appgo/model/common"
	"github.com/goecology/egoshop/appgo/model/constx"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/jinzhu/gorm"
)

const (
	goodsCommit = 1
	goodsCancel = 0
)

// 提交
func (u *userGoods) CreateOrUpdate(tx *gorm.DB, uid int, goods common.Goods, goodsType string) (err error) {
	var oneUserGoods mysql.UserGoods
	err = tx.Select("*").Where("uid = ? and goods_id =? and type_id = ?", uid, goods.Gid, goods.TypeId).Find(&oneUserGoods).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	// 创建该数据
	if err != nil && err == gorm.ErrRecordNotFound {
		oneUserGoods.Uid = uid
		oneUserGoods.GoodsId = goods.Gid
		oneUserGoods.TypeId = goods.TypeId
		// todo 优化，可以在创建就吧数据写入
		err = tx.Create(&oneUserGoods).Error
		if err != nil {
			return
		}
	}
	// 已经提交过，请报错误
	if isCommited(oneUserGoods, goodsType) {
		err = errors.New("already commited")
		return
	}

	switch goods.TypeId {
	case constx.TypeCom:
		err = tx.Table("com").Where("id = ?", goods.Gid).Updates(mysql.Ups{
			"cnt_" + goodsType: gorm.Expr("cnt_"+goodsType+"+?", 1),
		}).Error
		if err != nil {
			return
		}
	default:
		err = errors.New("type id is not exist")
		return
	}

	// 当为支付的时候，需要将预付款变为已付款
	if goodsType == "is_pay" {
		err = tx.Table("user_goods").Where("uid = ? and goods_id =?  and type_id = ?", uid, goods.Gid, goods.TypeId).Updates(mysql.Ups{
			goodsType:    goodsCommit,
			"is_pre_pay": 0,
		}).Error
	} else {
		err = tx.Table("user_goods").Where("uid = ? and goods_id =?  and type_id = ?", uid, goods.Gid, goods.TypeId).Updates(mysql.Ups{
			"is_" + goodsType: goodsCommit,
		}).Error
	}

	return
}

// 取消
func (u *userGoods) Cancel(tx *gorm.DB, uid int, goods common.Goods, goodsType string) (err error) {
	var oneUserGoods mysql.UserGoods
	err = tx.Select("*").Where("uid = ? and goods_id =? and type_id = ?", uid, goods.Gid, goods.TypeId).Find(&oneUserGoods).Error
	if err != nil {
		return
	}
	// 已经提交过，请报错误
	if isUnCommited(oneUserGoods, goodsType) {
		err = errors.New("uncommited")
		return
	}

	// todo 判断不好可能为负数
	switch goods.TypeId {
	case constx.TypeCom:
		err = tx.Table("com").Where("id = ?", goods.Gid).Updates(mysql.Ups{
			"cnt_" + goodsType: gorm.Expr("cnt_"+goodsType+"-?", 1),
		}).Error
		if err != nil {
			return
		}
	default:
		err = errors.New("type id is not exist")
		return
	}

	// 当为支付的时候，需要将预付款变为已付款
	if goodsType == "is_pay" {
		err = tx.Table("user_goods").Where("uid = ? and goods_id =?  and type_id = ?", uid, goods.Gid, goods.TypeId).Updates(mysql.Ups{
			goodsType:    goodsCancel,
			"is_pre_pay": 0,
		}).Error
	} else {
		err = tx.Table("user_goods").Where("uid = ? and goods_id =?  and type_id = ?", uid, goods.Gid, goods.TypeId).Updates(mysql.Ups{
			"is_" + goodsType: goodsCancel,
		}).Error
	}

	return
}

// 是否已经提交过
func isCommited(info mysql.UserGoods, goodsType string) bool {
	if goodsType == "star" && info.IsStar == 1 {
		return true
	}
	if goodsType == "collect" && info.IsCollect == 1 {
		return true
	}

	return false
}

// 是否已经提交过
func isUnCommited(info mysql.UserGoods, goodsType string) bool {
	if goodsType == "star" && info.IsStar == 0 {
		return true
	}
	if goodsType == "collect" && info.IsCollect == 0 {
		return true
	}

	return false
}
