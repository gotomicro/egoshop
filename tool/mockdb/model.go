package main

import "github.com/goecology/egoshop/appgo/model/mysql"

var Models []interface{}

func init() {
	Models = []interface{}{
		&mysql.Banner{},
		&mysql.Com{},
		&mysql.ComSku{},
		&mysql.ComStore{},
		&mysql.ComCate{},
		&mysql.ComRelateCate{},
		&mysql.Attachment{},
		&mysql.Cart{},
		&mysql.Address{},
		&mysql.AddressType{},
		&mysql.Order{},
		&mysql.OrderExtend{},
		&mysql.OrderLog{},
		&mysql.OrderPay{},
		&mysql.OrderGoods{},
		&mysql.ComImage{},
		&mysql.ComSpec{},
		&mysql.ComSpecValue{},
		&mysql.Freight{},
		&mysql.User{},
		&mysql.UserGoods{},
		&mysql.UserOpen{},
		&mysql.AccessToken{},
		&mysql.Comment{},
		&mysql.Banner{},
		&mysql.AdminUser{},
		&mysql.Signin{},
		&mysql.SigninLog{},
		&mysql.PointLog{},
		&mysql.PointLimit{},
		&mysql.Image{},
	}
}
